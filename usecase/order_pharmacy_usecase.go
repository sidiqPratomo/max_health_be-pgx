package usecase

import (
	"context"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type OrderPharmacyUsecase interface {
	GetOneOrderPharmacyById(ctx context.Context, orderPharmacyId int64) (*dto.OrderPharmacyResponse, error)
	GetAllOrderPharmacies(ctx context.Context, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error)
	GetAllUserOrderPharmacies(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error)
	GetAllPartnerOrderPharmacies(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error)
	GetAllPartnerOrderPharmaciesSummary(ctx context.Context, accountId int64) (*dto.AllOrderPharmaciesSummaryResponse, error)
	UpdateStatusToSent(ctx context.Context, accountId int64, orderPharmacyId int64) error
	UpdateStatusToConfirmed(ctx context.Context, accountId int64, orderPharmacyId int64) error
	UpdateStatusToCancelled(ctx context.Context, accountId int64, orderPharmacyId int64) error
}

type orderPharmacyUsecaseImpl struct {
	transaction               repository.Transaction
	orderPharmacyRepository   repository.OrderPharmacyRepository
	orderItemRepository       repository.OrderItemRepository
	userRepository            repository.UserRepository
	pharmacyManagerRepository repository.PharmacyManagerRepository
}

func NewOrderPharmacyUsecaseImpl(transaction repository.Transaction, orderPharmacyRepository repository.OrderPharmacyRepository, orderItemRepository repository.OrderItemRepository, userRepository repository.UserRepository, pharmacyManagerRepository repository.PharmacyManagerRepository) orderPharmacyUsecaseImpl {
	return orderPharmacyUsecaseImpl{
		transaction:               transaction,
		orderPharmacyRepository:   orderPharmacyRepository,
		orderItemRepository:       orderItemRepository,
		userRepository:            userRepository,
		pharmacyManagerRepository: pharmacyManagerRepository,
	}
}

func (u *orderPharmacyUsecaseImpl) GetOneOrderPharmacyById(ctx context.Context, orderPharmacyId int64) (*dto.OrderPharmacyResponse, error) {
	orderPharmacy, err := u.orderPharmacyRepository.FindOneById(ctx, orderPharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if orderPharmacy == nil {
		return nil, apperror.OrderNotFoundError()
	}

	orderItems, err := u.orderItemRepository.FindAllByOrderPharmacyId(ctx, orderPharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToOrderPharmacyResponse(*orderPharmacy, orderItems), err
}

func (u *orderPharmacyUsecaseImpl) GetAllOrderPharmacies(ctx context.Context, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error) {
	orderPharmacieIds, pageInfo, err := u.orderPharmacyRepository.FindAllIds(ctx, *validatedQuery)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	orderPharmaciesWithDetails := []*entity.OrderPharmacy{}

	if len(orderPharmacieIds) > 0 {
		orderPharmaciesWithDetails, err = u.orderPharmacyRepository.FindAllWithDetailsByIds(ctx, orderPharmacieIds)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
	}

	return dto.ConvertToAllOrderPharmaciesResponseWithPageInfoAndPointer(orderPharmaciesWithDetails, *pageInfo), nil
}

func (u *orderPharmacyUsecaseImpl) GetAllUserOrderPharmacies(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	orderPharmacies, pageInfo, err := u.orderPharmacyRepository.FindAllByOrderUserId(ctx, user.Id, *validatedQuery)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllOrderPharmaciesResponseWithPageInfo(orderPharmacies, *pageInfo), nil
}

func (u *orderPharmacyUsecaseImpl) GetAllPartnerOrderPharmacies(ctx context.Context, accountId int64, validatedQuery *util.ValidatedGetOrderQuery) (*dto.AllOrderPharmaciesResponse, error) {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return nil, apperror.PharmacyManagerNotFoundError()
	}

	orderPharmacyIds, pageInfo, err := u.orderPharmacyRepository.FindAllIdsByPharmacyManagerId(ctx, pharmacyManager.Id, *validatedQuery)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if len(orderPharmacyIds) == 0 {
		return dto.ConvertToAllOrderPharmaciesResponseWithPageInfoAndPointer([]*entity.OrderPharmacy{}, *pageInfo), nil
	}

	orderPharmaciesWithDetails, err := u.orderPharmacyRepository.FindAllWithDetailsByIds(ctx, orderPharmacyIds)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllOrderPharmaciesResponseWithPageInfoAndPointer(orderPharmaciesWithDetails, *pageInfo), nil
}

func (u *orderPharmacyUsecaseImpl) GetAllPartnerOrderPharmaciesSummary(ctx context.Context, accountId int64) (*dto.AllOrderPharmaciesSummaryResponse, error) {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return nil, apperror.PharmacyManagerNotFoundError()
	}

	orderPharmacySummary, err := u.orderPharmacyRepository.FindCountGroupedByOrderStatusIdByPharmacyManagerId(ctx, pharmacyManager.Id)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &dto.AllOrderPharmaciesSummaryResponse{
		AllCount:       orderPharmacySummary.AllCount,
		UnpaidCount:    orderPharmacySummary.UnpaidCount,
		ApprovalCount:  orderPharmacySummary.ApprovalCount,
		PendingCount:   orderPharmacySummary.PendingCount,
		SentCount:      orderPharmacySummary.SentCount,
		ConfirmedCount: orderPharmacySummary.ConfirmedCount,
		CanceledCount:  orderPharmacySummary.CanceledCount,
	}, nil
}

func (u *orderPharmacyUsecaseImpl) UpdateStatusToSent(ctx context.Context, accountId int64, orderPharmacyId int64) error {
	manager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if manager == nil {
		return apperror.PartnerNotFoundError()
	}

	orderPharmacy, err := u.orderPharmacyRepository.FindOneByOrderPharmacyId(ctx, orderPharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if orderPharmacy == nil {
		return apperror.PharmacyOrderNotFoundError()
	}

	orderManager, err := u.pharmacyManagerRepository.FindOneByPharmacyCourierId(ctx, orderPharmacy.PharmacyCourierId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if orderManager == nil {
		return apperror.PartnerNotFoundError()
	}

	if orderManager.Id != manager.Id {
		return apperror.ForbiddenAction()
	}

	if orderPharmacy.OrderStatusId != 3 {
		return apperror.InvalidOrderStatusError()
	}

	err = u.orderPharmacyRepository.UpdateOneStatusById(ctx, orderPharmacyId, 4)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *orderPharmacyUsecaseImpl) UpdateStatusToConfirmed(ctx context.Context, accountId int64, orderPharmacyId int64) error {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if user == nil {
		return apperror.UserNotFoundError()
	}

	orderPharmacy, err := u.orderPharmacyRepository.FindOneByOrderPharmacyId(ctx, orderPharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if orderPharmacy == nil {
		return apperror.PharmacyOrderNotFoundError()
	}

	if orderPharmacy.UserId != user.Id {
		return apperror.ForbiddenAction()
	}

	if orderPharmacy.OrderStatusId != 4 {
		return apperror.InvalidOrderStatusError()
	}

	err = u.orderPharmacyRepository.UpdateOneStatusById(ctx, orderPharmacyId, 5)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}

func (u *orderPharmacyUsecaseImpl) UpdateStatusToCancelled(ctx context.Context, accountId int64, orderPharmacyId int64) error {
	orderPharmacy, err := u.orderPharmacyRepository.FindOneByOrderPharmacyId(ctx, orderPharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if orderPharmacy == nil {
		return apperror.PharmacyOrderNotFoundError()
	}

	manager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if manager == nil {
		return apperror.PartnerNotFoundError()
	}

	orderManager, err := u.pharmacyManagerRepository.FindOneByPharmacyCourierId(ctx, orderPharmacy.PharmacyCourierId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if manager == nil {
		return apperror.PartnerNotFoundError()
	}

	if orderManager.Id != manager.Id {
		return apperror.ForbiddenAction()
	}

	if orderPharmacy.OrderStatusId >= 6 || orderPharmacy.OrderStatusId < 3 {
		return apperror.InvalidOrderStatusError()
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	pharmacyDrugRepo := tx.PharmacyDrugRepo()
	orderPharmacyRepo := tx.OrderPharmacyRepository()
	stockChangeRepo := tx.StockChangeRepo()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = orderPharmacyRepo.UpdateOneStatusById(ctx, orderPharmacyId, 6)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	stockChanges, err := pharmacyDrugRepo.UpdatePharmacyDrugsByOrderPharmacyId(ctx, orderPharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = stockChangeRepo.PostStockChanges(ctx, stockChanges)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
