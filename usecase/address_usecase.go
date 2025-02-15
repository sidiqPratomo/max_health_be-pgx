package usecase

import (
	"context"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type AddressUsecase interface {
	GetAllProvinces(ctx context.Context) (*dto.AllProvincesResponse, error)
	GetAllCitiesByProvinceCode(ctx context.Context, provinceCode string) (*dto.AllCitiesResponse, error)
	GetAllDistrictByCityCode(ctx context.Context, cityCode string) (*dto.AllDistrictsResponse, error)
	GetAllSubdistrictByDistrictCode(ctx context.Context, districtCode string) (*dto.AllSubdistrictsResponse, error)
}

type addressUsecaseImpl struct {
	AddressRepository repository.AddressRepository
}

func NewAddressUsecaseImpl(addressRepository repository.AddressRepository) addressUsecaseImpl {
	return addressUsecaseImpl{
		AddressRepository: addressRepository,
	}
}

func (u *addressUsecaseImpl) GetAllProvinces(ctx context.Context) (*dto.AllProvincesResponse, error) {
	provinces, err := u.AddressRepository.FindAllProvinces(ctx)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllProvincesResponse(provinces), nil
}

func (u *addressUsecaseImpl) GetAllCitiesByProvinceCode(ctx context.Context, provinceCode string) (*dto.AllCitiesResponse, error) {
	cities, err := u.AddressRepository.FindAllCitiesByProvinceCode(ctx, provinceCode)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllCitiesResponse(cities), nil
}

func (u *addressUsecaseImpl) GetAllDistrictByCityCode(ctx context.Context, cityCode string) (*dto.AllDistrictsResponse, error) {
	districts, err := u.AddressRepository.FindAllDistrictsByCityCode(ctx, cityCode)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllDistrictsResponse(districts), nil
}

func (u *addressUsecaseImpl) GetAllSubdistrictByDistrictCode(ctx context.Context, districtCode string) (*dto.AllSubdistrictsResponse, error) {
	subdistricts, err := u.AddressRepository.FindAllSubdistrictsByDistrictCode(ctx, districtCode)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return dto.ConvertToAllSubdistrictsResponse(subdistricts), nil
}
