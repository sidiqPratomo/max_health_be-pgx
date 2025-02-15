package usecase

import (
	"context"

	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type DrugFormUsecase interface {
	GetAllDrugForm(ctx context.Context) ([]dto.DrugForm, error)
}

type drugFormUsecaseImpl struct {
	drugFormRepository repository.DrugFormRepository
}

func NewdrugFormUsecaseImpl(drugFormRepository repository.DrugFormRepository) drugFormUsecaseImpl {
	return drugFormUsecaseImpl{
		drugFormRepository: drugFormRepository,
	}
}

func (u *drugFormUsecaseImpl) GetAllDrugForm(ctx context.Context) ([]dto.DrugForm, error) {
	drugFormList, err := u.drugFormRepository.GetAllDrugForm(ctx)
	if err != nil {
		return nil, err
	}

	drugFormListDTO := dto.ConvertToDrugFormListDTO(drugFormList)

	return drugFormListDTO, nil
}
