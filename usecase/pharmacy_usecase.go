package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type PharmacyUsecase interface {
	GetAllPharmacyByManagerId(ctx context.Context, accountId int64, limit string, page string, search string) (*dto.GetAllPharmacyResponse, error)
	CreateOnePharmacy(ctx context.Context, accountId int64, pharmacyRequest dto.PharmacyRequest) error
	UpdateOnePharmacy(ctx context.Context, accountId int64, updatePharmacyRequest dto.UpdatePharmacyRequest) error
	DeleteOnePharmacyById(ctx context.Context, accountId int64, pharmacyId int64) error
	AdminGetAllPharmacyByManagerId(ctx context.Context, managerId int64, limit string, page string, search string) (*dto.GetAllPharmacyResponse, error)
}

type pharmacyUsecaseImpl struct {
	pharmacyManagerRepository repository.PharmacyManagerRepository
	pharmacyRepository        repository.PharmacyRepository
	pharmacyDrugRepository    repository.PharmacyDrugRepository
	addressRepository         repository.AddressRepository
	courierRepository         repository.CourierRepository
	orderPharmacyRepository   repository.OrderPharmacyRepository
	transaction               repository.Transaction
}

func NewPharmacyUsecaseImpl(pharmacyManagerRepository repository.PharmacyManagerRepository, pharmacyRepository repository.PharmacyRepository, pharmacyDrugRepository repository.PharmacyDrugRepository, addressRepository repository.AddressRepository, courierRepository repository.CourierRepository, orderPharmacyRepository repository.OrderPharmacyRepository, transaction repository.Transaction) pharmacyUsecaseImpl {
	return pharmacyUsecaseImpl{
		pharmacyManagerRepository: pharmacyManagerRepository,
		pharmacyRepository:        pharmacyRepository,
		pharmacyDrugRepository:    pharmacyDrugRepository,
		addressRepository:         addressRepository,
		courierRepository:         courierRepository,
		orderPharmacyRepository:   orderPharmacyRepository,
		transaction:               transaction,
	}
}

func (u *pharmacyUsecaseImpl) GetAllPharmacyByManagerId(ctx context.Context, accountId int64, limit string, page string, search string) (*dto.GetAllPharmacyResponse, error) {
	manager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if manager == nil {
		return nil, apperror.PartnerNotFoundError()
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

	pharmacies, pageInfo, err := u.pharmacyRepository.FindAllByManagerId(ctx, manager.Id, limitInt, offset, search)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pharmaciesDto := []dto.PharmacyDTO{}
	for _, pharmacy := range pharmacies {
		pharmacyDTO := dto.PharmacyDTO{
			Id:                      pharmacy.Id,
			PharmacyManagerId:       pharmacy.PharmacyManagerId,
			Name:                    pharmacy.Name,
			PharmacistName:          pharmacy.PharmacistName,
			PharmacistLicenseNumber: pharmacy.PharmacistLicenseNumber,
			PharmacistPhoneNumber:   pharmacy.PharmacistPhoneNumber,
			City:                    pharmacy.City,
			Address:                 pharmacy.Address,
			Latitude:                pharmacy.Latitude,
			Longitude:               pharmacy.Longitude,
		}
		pharmaciesDto = append(pharmaciesDto, pharmacyDTO)
	}

	pharmacyResponse := dto.GetAllPharmacyResponse{
		PageInfo:   *pageInfo,
		Pharmacies: pharmaciesDto,
	}

	return &pharmacyResponse, nil
}

func (u *pharmacyUsecaseImpl) CreateOnePharmacy(ctx context.Context, accountId int64, pharmacyRequest dto.PharmacyRequest) error {
	pharmacy := dto.ConvertPharmacyRequestToPharmacy(pharmacyRequest)

	pharmacy.Name = strings.Trim(pharmacy.Name, " ")

	pharmacyManager, err := u.pharmacyManagerRepository.FindOneById(ctx, pharmacyRequest.PharmacyManagerId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return apperror.PharmacyManagerNotFoundError()
	}

	couriers, err := u.courierRepository.FindAll(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	courierIds := []int64{}
	for _, courier := range couriers {
		courierIds = append(courierIds, courier.Id)
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacy.PharmacyManagerId = pharmacyManager.Id

	pharmacyRepo := tx.PharmacyRepository()
	pharmacyOperationalRepo := tx.PharmacyOperationalRepository()
	pharmacyCourierRepo := tx.PharmacyCourierRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	pharmacyId, err := pharmacyRepo.CreateOne(ctx, &pharmacy)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	if err = pharmacyOperationalRepo.CreateBulk(ctx, *pharmacyId, days); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = pharmacyCourierRepo.CreateBulk(ctx, *pharmacyId, courierIds); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *pharmacyUsecaseImpl) UpdateOnePharmacy(ctx context.Context, accountId int64, updatePharmacyRequest dto.UpdatePharmacyRequest) error {
	newPharmacy := dto.ConvertUpdatePharmacyRequestToPharmacy(updatePharmacyRequest)

	if newPharmacy.Id < 1 {
		return apperror.PharmacyNotFoundError()
	}

	newPharmacy.Name = strings.Trim(newPharmacy.Name, " ")
	newPharmacy.PharmacistName = strings.Trim(newPharmacy.PharmacistName, " ")
	newPharmacy.PharmacistLicenseNumber = strings.Trim(newPharmacy.PharmacistLicenseNumber, " ")
	newPharmacy.PharmacistPhoneNumber = strings.Trim(newPharmacy.PharmacistPhoneNumber, " ")
	newPharmacy.Address = strings.Trim(newPharmacy.Address, " ")
	newPharmacy.City = strings.Trim(newPharmacy.City, " ")

	if len(updatePharmacyRequest.Operationals) != 7 {
		return apperror.InvalidPharmacyOperationalError()
	}

	if len(updatePharmacyRequest.Couriers) == 0 {
		return apperror.InvalidPharmacyCourierError()
	}

	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return apperror.PharmacyManagerNotFoundError()
	}

	oldPharmacy, err := u.pharmacyRepository.FindOneById(ctx, newPharmacy.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if oldPharmacy == nil {
		return apperror.PharmacyNotFoundError()
	}

	if pharmacyManager.Id != oldPharmacy.PharmacyManagerId {
		return apperror.ForbiddenAction()
	}

	if newPharmacy.City != oldPharmacy.City && newPharmacy.City != "" {
		newPharmacy.City = strings.ToLower(newPharmacy.City)
		if strings.Contains(newPharmacy.City, "kota") || strings.Contains(newPharmacy.City, "kabupaten") {
			newPharmacy.City = strings.ReplaceAll(newPharmacy.City, "kota", "")
			newPharmacy.City = strings.ReplaceAll(newPharmacy.City, "kabupaten", "")
		}
		cityId, err := u.addressRepository.FindOneCityByName(ctx, "%"+newPharmacy.City)
		if err != nil {
			return apperror.InternalServerError(err)
		}
		if cityId == nil {
			return apperror.LocationNotFoundError()
		}
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacyRepo := tx.PharmacyRepository()
	pharmacyOperationalRepo := tx.PharmacyOperationalRepository()
	pharmacyCourierRepo := tx.PharmacyCourierRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err = pharmacyRepo.UpdateOne(ctx, newPharmacy); err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacyOperationals := dto.AllUpdatePharmacyOperationalsToAllPharmacyOperationals(updatePharmacyRequest.Operationals)

	for i := 0; i < len(pharmacyOperationals); i++ {
		if err = pharmacyOperationalRepo.UpdateOneById(ctx, pharmacyOperationals[i]); err != nil {
			return apperror.InternalServerError(err)
		}
	}

	couriers := dto.AllUpdatePharmacyCourierRequestToAllPharmacyCouriers(updatePharmacyRequest.Couriers)

	for i := 0; i < len(couriers); i++ {
		if err = pharmacyCourierRepo.UpdateOneById(ctx, couriers[i]); err != nil {
			return apperror.InternalServerError(err)
		}
	}

	return nil
}

func (u *pharmacyUsecaseImpl) DeleteOnePharmacyById(ctx context.Context, accountId int64, pharmacyId int64) error {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return apperror.PharmacyManagerNotFoundError()
	}

	oldPharmacy, err := u.pharmacyRepository.FindOneById(ctx, pharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if oldPharmacy == nil {
		return apperror.PharmacyNotFoundError()
	}
	if pharmacyManager.Id != oldPharmacy.PharmacyManagerId {
		return apperror.ForbiddenAction()
	}

	ongoingOrderPharmacyIds, err := u.orderPharmacyRepository.FindAllOngoingIdsByPharmacyId(ctx, pharmacyId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if len(ongoingOrderPharmacyIds) > 0 {
		return apperror.OngoingOrderExistsError()
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	pharmacyRepo := tx.PharmacyRepository()
	pharmacyOperationalRepo := tx.PharmacyOperationalRepository()
	pharmacyCourierRepo := tx.PharmacyCourierRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err = pharmacyRepo.DeleteOneById(ctx, pharmacyId); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = pharmacyOperationalRepo.DeleteBulkByPharmacyId(ctx, pharmacyId); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = pharmacyCourierRepo.DeleteBulkByPharmacyId(ctx, pharmacyId); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *pharmacyUsecaseImpl) AdminGetAllPharmacyByManagerId(ctx context.Context, managerId int64, limit string, page string, search string) (*dto.GetAllPharmacyResponse, error) {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneById(ctx, managerId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if pharmacyManager == nil {
		return nil, apperror.UnauthorizedError()
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

	pharmacies, pageInfo, err := u.pharmacyRepository.FindAllByManagerId(ctx, managerId, limitInt, offset, search)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	pharmaciesDto := []dto.PharmacyDTO{}
	for _, pharmacy := range pharmacies {
		pharmacyDTO := dto.PharmacyDTO{
			Id:                      pharmacy.Id,
			PharmacyManagerId:       pharmacy.PharmacyManagerId,
			Name:                    pharmacy.Name,
			PharmacistName:          pharmacy.PharmacistName,
			PharmacistLicenseNumber: pharmacy.PharmacistLicenseNumber,
			PharmacistPhoneNumber:   pharmacy.PharmacistPhoneNumber,
			City:                    pharmacy.City,
			Address:                 pharmacy.Address,
			Latitude:                pharmacy.Latitude,
			Longitude:               pharmacy.Longitude,
		}
		pharmaciesDto = append(pharmaciesDto, pharmacyDTO)
	}

	pharmacyResponse := dto.GetAllPharmacyResponse{
		PageInfo:   *pageInfo,
		Pharmacies: pharmaciesDto,
	}

	return &pharmacyResponse, nil
}
