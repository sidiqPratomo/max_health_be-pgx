package usecase

import (
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type OrderUsecase interface {
	CheckoutOrder(ctx context.Context, orderCheckoutRequest dto.OrderCheckoutRequest) (*int64, error)
	ConfirmPayment(ctx context.Context, orderId int64, statusId int64) error
	UploadPaymentProofOrder(ctx context.Context, accountId int64, orderId int64, file multipart.File, fileHeader multipart.FileHeader) error
	GetAllUserPendingOrders(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrdersResponse, error)
	GetAllOrders(ctx context.Context, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrdersResponse, error)
	GetOneOrderById(ctx context.Context, orderId int64) (*dto.OrderResponse, error)
	CancelOrder(ctx context.Context, accountId int64, orderId int64) error
}

type orderUsecaseImpl struct {
	transaction             repository.Transaction
	userRepository          repository.UserRepository
	orderRepository         repository.OrderRepository
	orderPharmacyRepository repository.OrderPharmacyRepository
}

func NewOrderUsecaseImpl(transaction repository.Transaction, userRepository repository.UserRepository, orderRepository repository.OrderRepository, orderPharmacyRepository repository.OrderPharmacyRepository) orderUsecaseImpl {
	return orderUsecaseImpl{
		transaction:             transaction,
		userRepository:          userRepository,
		orderRepository:         orderRepository,
		orderPharmacyRepository: orderPharmacyRepository,
	}
}

func (u *orderUsecaseImpl) CheckoutOrder(ctx context.Context, orderCheckoutRequest dto.OrderCheckoutRequest) (*int64, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, orderCheckoutRequest.AccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	for _, pharmacy := range orderCheckoutRequest.Pharmacies {
		if len(pharmacy.CartItemIds) < 1 {
			return nil, apperror.EmptyCartSelectionError()
		}
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	orderRepo := tx.OrderRepository()
	orderPharmacyRepo := tx.OrderPharmacyRepository()
	orderItemRepo := tx.OrderItemRepository()
	cartRepo := tx.CartRepository()
	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	stockChangeRepo := tx.StockChangeRepo()
	stockMutationRepo := tx.StockMutationRepo()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	orderId, err := orderRepo.PostOneOrder(ctx, user.Id, orderCheckoutRequest.Address, orderCheckoutRequest.TotalAmount)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	orderPharmacies, err := orderPharmacyRepo.PostOrderPharmacies(ctx, orderId, orderCheckoutRequest)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	for i, pharmacy := range orderCheckoutRequest.Pharmacies {
		orderPharmacies[i].CartItems, err = cartRepo.GetAllCartDetailByIds(ctx, pharmacy.CartItemIds)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	err = orderItemRepo.PostOrderItems(ctx, orderPharmacies)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	allCartItems := []entity.CartItemForCheckout{}
	for _, pharmacy := range orderPharmacies {
		allCartItems = append(allCartItems, pharmacy.CartItems...)
	}

	err = pharmacyDrugRepo.GetPharmacyDrugsByCartForUpdate(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	carts, err := cartRepo.GetAllCartsForChangesByCartIds(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	err = stockChangeRepo.PostStockChangesByCartIds(ctx, carts)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pharmacyDrugs, err := pharmacyDrugRepo.UpdatePharmacyDrugsByCartId(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	alternatives, err := stockMutationRepo.GetPossibleStockMutation(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	stockMutationList := []entity.PossibleStockMutation{}
	stockChangesList := []entity.StockChange{}
	insufficientCartItems := []int64{}
	for _, pharmacyDrug := range pharmacyDrugs {
		if pharmacyDrug.Stock >= 0 {
			continue
		}
		stock := pharmacyDrug.Stock
		for _, alternative := range alternatives {
			if alternative.CartItemId != pharmacyDrug.CartId {
				continue
			}
			if stock+alternative.AlternativeStock < 0 {
				stock += alternative.AlternativeStock
				stockMutationList = append(stockMutationList, alternative)
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.OriginalPharmacyDrug,
					FinalStock: stock, Amount: alternative.AlternativeStock})
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.AlternativePharmacy,
					FinalStock: 0, Amount: -1 * alternative.AlternativeStock})
			} else {
				partialAlternative := alternative
				partialAlternative.AlternativeStock = stock * -1
				stockMutationList = append(stockMutationList, partialAlternative)
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.OriginalPharmacyDrug,
					FinalStock: 0, Amount: partialAlternative.AlternativeStock})
				stockChangesList = append(stockChangesList, entity.StockChange{PharmacyDrugId: alternative.AlternativePharmacy,
					FinalStock: alternative.AlternativeStock - partialAlternative.AlternativeStock,
					Amount:     -1 * partialAlternative.AlternativeStock})
				stock = 0
				break
			}
		}
		if stock < 0 {
			insufficientCartItems = append(insufficientCartItems, pharmacyDrug.CartId)
		}
	}

	if len(insufficientCartItems) > 0 {
		err = apperror.InsufficientStockDuringCheckoutError(insufficientCartItems)
		return nil, err
	}

	if len(stockMutationList) > 0 {
		err = stockMutationRepo.PostStockMutations(ctx, stockMutationList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		err = stockChangeRepo.PostStockChangesFromMutation(ctx, stockChangesList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}

		err = pharmacyDrugRepo.UpdatePharmacyDrugsForStockMutation(ctx, stockChangesList)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	err = cartRepo.DeleteCarts(ctx, allCartItems)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &orderId, nil
}

func (u *orderUsecaseImpl) ConfirmPayment(ctx context.Context, orderId int64, statusId int64) error {
	orderPharmacies, err := u.orderPharmacyRepository.FindAllByOrderId(ctx, orderId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if len(orderPharmacies) == 0 {
		return apperror.OrderNotFoundError()
	}

	for i := 0; i < len(orderPharmacies); i++ {
		if orderPharmacies[i].OrderStatusId != 2 {
			return apperror.InvalidOrderStatusError()
		}
	}

	order, err := u.orderRepository.FindOneOrderByOrderId(ctx, orderId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if order.PaymentProof == "" {
		return apperror.PaymentProofIsEmptyError()
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	orderRepo := tx.OrderRepository()
	orderPharmacyRepo := tx.OrderPharmacyRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if statusId == 1 {
		if strings.Split(order.PaymentProof, "/")[2] == "res.cloudinary.com" {
			util.DeleteInCloudinary(order.PaymentProof)
		}
		if err := orderRepo.UpdatePaymentProofOne(ctx, &entity.Order{
			Id:           orderId,
			PaymentProof: "",
		}); err != nil {
			return apperror.InternalServerError(err)
		}
	}

	if err := orderPharmacyRepo.UpdateStatusBulkByOrderId(ctx, orderId, statusId); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *orderUsecaseImpl) UploadPaymentProofOrder(ctx context.Context, accountId int64, orderId int64, file multipart.File, fileHeader multipart.FileHeader) error {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if user == nil {
		return apperror.UserNotFoundError()
	}

	orderPharmacies, err := u.orderPharmacyRepository.FindAllByOrderId(ctx, orderId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if len(orderPharmacies) == 0 {
		return apperror.OrderNotFoundError()
	}
	if orderPharmacies[0].UserId != user.Id {
		return apperror.ForbiddenAction()
	}

	for i := 0; i < len(orderPharmacies); i++ {
		if orderPharmacies[i].OrderStatusId != 1 {
			return apperror.InvalidOrderStatusError()
		}
	}

	filePath, _, err := util.ValidateFile(fileHeader, appconstant.OrderPaymentProofsUrl, []string{"png", "jpg", "jpeg"}, 2000000)
	if err != nil {
		return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
	}

	paymentProofUrl, err := util.UploadToCloudinary(file, *filePath)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	orderRepo := tx.OrderRepository()
	orderPharmacyRepo := tx.OrderPharmacyRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err := orderRepo.UpdatePaymentProofOne(ctx, &entity.Order{
		Id:           orderId,
		PaymentProof: paymentProofUrl,
	}); err != nil {
		return apperror.InternalServerError(err)
	}

	if err := orderPharmacyRepo.UpdateStatusBulkByOrderId(ctx, orderId, 2); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *orderUsecaseImpl) GetOneOrderById(ctx context.Context, orderId int64) (*dto.OrderResponse, error) {
	order, err := u.orderRepository.FindOneOrderByOrderId(ctx, orderId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if order == nil {
		return nil, apperror.OrderNotFoundError()
	}

	res := dto.ConvertToOrderResponse(*order)

	return &res, err
}

func (u *orderUsecaseImpl) GetAllOrders(ctx context.Context, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrdersResponse, error) {
	orderIds, pageInfo, err := u.orderRepository.FindAll(ctx, *validatedQuery)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	ordersWithDetails := []*entity.Order{}

	if len(orderIds) > 0 {
		ordersWithDetails, err = u.orderRepository.FindAllWithDetails(ctx, orderIds)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	return dto.ConvertToAllOrdersResponse(ordersWithDetails, *pageInfo), nil
}

func (u *orderUsecaseImpl) GetAllUserPendingOrders(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrdersResponse, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	orderIds, pageInfo, err := u.orderRepository.FindAllPendingByUserId(ctx, user.Id, *validatedQuery)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	ordersWithDetails, err := u.orderRepository.FindAllPendingWithDetailsByUserId(ctx, user.Id, orderIds)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllOrdersResponse(ordersWithDetails, *pageInfo), nil
}

func (u *orderUsecaseImpl) CancelOrder(ctx context.Context, accountId int64, orderId int64) error {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if user == nil {
		return apperror.UserNotFoundError()
	}

	orderPharmacies, err := u.orderPharmacyRepository.FindAllByOrderId(ctx, orderId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if len(orderPharmacies) == 0 {
		return apperror.OrderNotFoundError()
	}
	if orderPharmacies[0].UserId != user.Id {
		return apperror.ForbiddenAction()
	}

	for i := 0; i < len(orderPharmacies); i++ {
		if orderPharmacies[i].OrderStatusId != 1 {
			return apperror.InvalidOrderStatusError()
		}
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	orderPharmacyRepo := tx.OrderPharmacyRepository()
	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	stockChangeRepo := tx.StockChangeRepo()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()
	if err := orderPharmacyRepo.UpdateStatusBulkByOrderId(ctx, orderId, 6); err != nil {
		return apperror.InternalServerError(err)
	}
	stockChanges, err := pharmacyDrugRepo.UpdatePharmacyDrugsByOrderId(ctx, orderId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = stockChangeRepo.PostStockChanges(ctx, stockChanges)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
