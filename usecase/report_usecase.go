package usecase

import (
	"context"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type ReportUsecase interface {
	GetPharmacyDrugCategoryReport(ctx context.Context, accountId int64, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugCategorySalesVolumeRevenueResponse, error)
	GetPharmacyDrugReport(ctx context.Context, accountId int64, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugSalesVolumeRevenueResponse, error)
	GetDrugCategoryReport(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugCategorySalesVolumeRevenueResponse, error)
	GetDrugReport(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugSalesVolumeRevenueResponse, error)
}

type reportUsecaseImpl struct {
	orderItemRepository       repository.OrderItemRepository
	pharmacyRepository        repository.PharmacyRepository
	pharmacyManagerRepository repository.PharmacyManagerRepository
}

func NewreportUsecaseImpl(orderItemRepository repository.OrderItemRepository, pharmacyRepository repository.PharmacyRepository, pharmacyManagerRepository repository.PharmacyManagerRepository) reportUsecaseImpl {
	return reportUsecaseImpl{
		orderItemRepository:       orderItemRepository,
		pharmacyRepository:        pharmacyRepository,
		pharmacyManagerRepository: pharmacyManagerRepository,
	}
}

func (u *reportUsecaseImpl) GetPharmacyDrugCategoryReport(ctx context.Context, accountId int64, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugCategorySalesVolumeRevenueResponse, error) {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return nil, apperror.PharmacyManagerNotFoundError()
	}

	pharmacy, err := u.pharmacyRepository.GetOnePharmacyByPharmacyId(ctx, validatedGetReportQuery.PharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if pharmacy == nil {
		return nil, apperror.PharmacyNotFoundError()
	}

	if pharmacy.PharmacyManagerId != pharmacyManager.Id {
		return nil, apperror.ForbiddenAction()
	}

	pharmacyDrugCategoryReport, err := u.orderItemRepository.FindPharmacyDrugCategorySalesVolumeRevenueByPharmacyId(ctx, validatedGetReportQuery)
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllDrugCategorySalesVolumeRevenueResponse(pharmacyDrugCategoryReport), err
}

func (u *reportUsecaseImpl) GetPharmacyDrugReport(ctx context.Context, accountId int64, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugSalesVolumeRevenueResponse, error) {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return nil, apperror.PharmacyManagerNotFoundError()
	}

	pharmacy, err := u.pharmacyRepository.GetOnePharmacyByPharmacyId(ctx, validatedGetReportQuery.PharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if pharmacy == nil {
		return nil, apperror.PharmacyNotFoundError()
	}

	if pharmacy.PharmacyManagerId != pharmacyManager.Id {
		return nil, apperror.ForbiddenAction()
	}

	pharmacyDrugReport, err := u.orderItemRepository.FindPharmacyDrugSalesVolumeRevenueByPharmacyId(ctx, validatedGetReportQuery)
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllDrugSalesVolumeRevenueResponse(pharmacyDrugReport), err
}

func (u *reportUsecaseImpl) GetDrugCategoryReport(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugCategorySalesVolumeRevenueResponse, error) {
	pharmacy, err := u.pharmacyRepository.GetOnePharmacyByPharmacyId(ctx, validatedGetReportQuery.PharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacy == nil {
		return nil, apperror.PharmacyNotFoundError()
	}

	drugCategoryReport, err := u.orderItemRepository.FindPharmacyDrugCategorySalesVolumeRevenueByPharmacyId(ctx, validatedGetReportQuery)
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllDrugCategorySalesVolumeRevenueResponse(drugCategoryReport), err
}

func (u *reportUsecaseImpl) GetDrugReport(ctx context.Context, validatedGetReportQuery util.ValidatedGetReportQuery) (*dto.AllDrugSalesVolumeRevenueResponse, error) {
	pharmacy, err := u.pharmacyRepository.GetOnePharmacyByPharmacyId(ctx, validatedGetReportQuery.PharmacyId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if pharmacy == nil {
		return nil, apperror.PharmacyNotFoundError()
	}

	drugReport, err := u.orderItemRepository.FindPharmacyDrugSalesVolumeRevenueByPharmacyId(ctx, validatedGetReportQuery)
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllDrugSalesVolumeRevenueResponse(drugReport), err
}
