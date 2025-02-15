package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type CartUsecase interface {
	CalculateDeliveryFee(ctx context.Context, deliveryFeeRequest dto.DeliveryFeeRequest) (*dto.AllDeliveryFeeResponse, error)
	CreateOneCart(ctx context.Context, pharmacyDrugId int64) error
	UpdateOneCart(ctx context.Context, cartItemID int64, quantity int) error
	DeleteOneCart(ctx context.Context, cartItemID int64) error
	GetAllCartById(ctx context.Context, page string, limit string) (*dto.CartDTOResponse, error)
}

type cartUsecaseImpl struct {
	userRepository         repository.UserRepository
	userAddressRepository  repository.UserAddressRepository
	cartRepository         repository.CartRepository
	pharmacyDrugRepository repository.PharmacyDrugRepository
}

func NewCartUsecaseImpl(pharmacyDrugRepository repository.PharmacyDrugRepository, userRepository repository.UserRepository, userAddressRepository repository.UserAddressRepository, cartRepository repository.CartRepository) cartUsecaseImpl {
	return cartUsecaseImpl{
		userRepository:         userRepository,
		userAddressRepository:  userAddressRepository,
		cartRepository:         cartRepository,
		pharmacyDrugRepository: pharmacyDrugRepository,
	}
}

func (u *cartUsecaseImpl) CalculateDeliveryFee(ctx context.Context, deliveryFeeRequest dto.DeliveryFeeRequest) (*dto.AllDeliveryFeeResponse, error) {
	if len(deliveryFeeRequest.CartItemsId) == 0 {
		return nil, apperror.EmptyCartSelectionError()
	}

	user, err := u.userRepository.FindUserByAccountId(ctx, deliveryFeeRequest.AccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if user == nil {
		return nil, apperror.AccountNotFoundError()
	}

	address, err := u.userAddressRepository.GetOneUserAddressByAddressId(ctx, deliveryFeeRequest.UserAddressId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if address == nil {
		return nil, apperror.UserAddressNotFoundError()
	}

	carts, err := u.cartRepository.GetCartsByIds(ctx, deliveryFeeRequest.CartItemsId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if len(carts) != len(deliveryFeeRequest.CartItemsId) {
		return nil, apperror.CartItemNotFoundError()
	}

	for _, cart := range carts {
		if cart.UserId != user.Id {
			return nil, apperror.UnauthorizedUserCartAccessError()
		}
	}

	deliveryFees, err := u.cartRepository.GetPharmacyDeliveryFeeForCart(ctx, deliveryFeeRequest.CartItemsId, deliveryFeeRequest.UserAddressId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	deliveryFeesResponse := dto.AllDeliveryFeeResponse{Pharmacies: deliveryFees}

	return &deliveryFeesResponse, nil
}

func (u *cartUsecaseImpl) CreateOneCart(ctx context.Context, pharmacyDrugId int64) error {
	id := appconstant.AccountId

	accountId := ctx.Value(id)
	if accountId == nil {
		return errors.New("userId is nil")
	}

	accountID, ok := accountId.(int64)
	if !ok {
		return errors.New("userId is not of type int in the context")
	}

	if pharmacyDrugId < 1 {
		return errors.New("pharmacy_drug_id can't be less than 1")
	}

	pharmacyDrug, err := u.pharmacyDrugRepository.GetPharmacyDrugById(ctx, pharmacyDrugId)
	if err != nil {
		return apperror.BadRequestError(err)
	}

	if pharmacyDrug.Stock < 1 {
		return apperror.NewAppError(422, errors.New("insufficient stock"), "insufficient stock")
	}

	_, err = u.cartRepository.PostOneCart(ctx, accountID, pharmacyDrugId, 1)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *cartUsecaseImpl) UpdateOneCart(ctx context.Context, cartItemID int64, quantity int) error {
	id := appconstant.AccountId

	accountId := ctx.Value(id)
	if accountId == nil {
		return errors.New("userId is nil")
	}

	accountID, ok := accountId.(int64)
	if !ok {
		return errors.New("userId is not of type int in the context")
	}

	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	stock, err := u.cartRepository.GetStockByCartId(ctx, cartItemID)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if quantity > *stock {
		return apperror.NewAppError(422, errors.New("insufficient stock"), "insufficient stock")
	}

	err = u.cartRepository.UpdateOneCart(ctx, accountID, cartItemID, quantity)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *cartUsecaseImpl) DeleteOneCart(ctx context.Context, cartItemID int64) error {
	id := appconstant.AccountId

	accountId := ctx.Value(id)
	if accountId == nil {
		return errors.New("userId is nil")
	}

	accountID, ok := accountId.(int64)
	if !ok {
		return errors.New("userId is not of type int in the context")
	}

	err := u.cartRepository.DeleteOneCart(ctx, accountID, cartItemID)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *cartUsecaseImpl) GetAllCartById(ctx context.Context, limit string, page string) (*dto.CartDTOResponse, error) {
	id := appconstant.AccountId
	accountId := ctx.Value(id)
	if accountId == nil {
		return nil, errors.New("userId is nil")
	}

	accountID, ok := accountId.(int64)
	if !ok {
		return nil, errors.New("userId is not of type int in the context")
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 12
	}

	offset := (pageInt - 1) * limitInt
	if offset < 0 {
		offset = 0
	}

	carts, pageInfo, err := u.cartRepository.GetAllCart(ctx, accountID, limit, offset)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	GetAllCart := []dto.CartDTO{}

	for _, cart := range carts {

		cartDto := dto.CartDTO{
			Id:       cart.CartItemId,
			UserId:   cart.UserId,
			Quantity: cart.Quantity,
			PharmacyDrugs: dto.PharmacyDrugsDto{
				Id:           cart.PharmacyDrugId,
				PharmacyId:   cart.PharmacyId,
				Name:         cart.DrugName,
				Price:        cart.Price,
				Image:        cart.Image,
				PharmacyName: cart.PharmacyName,
				Stock:        cart.Stock,
			},
		}

		GetAllCart = append(GetAllCart, cartDto)
	}
	cartResponse := dto.CartDTOResponse{
		Page:  *pageInfo,
		Carts: GetAllCart,
	}

	return &cartResponse, nil
}
