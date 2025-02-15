package usecase

import (
	"context"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type StockUsecase interface {
	GetAllStockChanges(ctx context.Context, accountId int64, pharmacyId *int64) ([]dto.StockChangeResponse, error)
}

type stockUsecaseImpl struct {
	stockRepository   repository.StockChangeRepository
	managerRepository repository.PharmacyManagerRepository
}

func NewStockUsecaseImpl(stockRepository repository.StockChangeRepository, managerRepository repository.PharmacyManagerRepository) stockUsecaseImpl {
	return stockUsecaseImpl{
		stockRepository:   stockRepository,
		managerRepository: managerRepository,
	}
}

func (u *stockUsecaseImpl) GetAllStockChanges(ctx context.Context, accountId int64, pharmacyId *int64) ([]dto.StockChangeResponse, error) {
	manager, err := u.managerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if manager == nil {
		return nil, apperror.PharmacyManagerNotFoundError()
	}

	stockChanges, err := u.stockRepository.GetStockChanges(ctx, manager.Id, pharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	return stockChanges, nil
}
