package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
)

type UserAddressUsecase interface {
	AddUserAddress(ctx context.Context, accountId int64, addUserAddressRequest dto.AddUserAddressRequest) error
	AddUserAddressAutofill(ctx context.Context, accountId int64, addUserAddressAutofillRequest dto.AddUserAddressAutofillRequest) error
	DeleteUserAddress(ctx context.Context, userAddressId int64, accountId int64) error
	GetAllUserAddress(ctx context.Context, accountId int64) (*dto.AllUserAddressResponse, error)
	UpdateUserAddress(ctx context.Context, accountId int64, updateUserAddressRequest dto.UpdateUserAddressRequest) error
}

type userAddressUsecaseImpl struct {
	userRepository        repository.UserRepository
	userAddressRepository repository.UserAddressRepository
	addressRepository     repository.AddressRepository
	transaction           repository.Transaction
}

func NewUserAddressUsecaseImpl(userRepository repository.UserRepository, userAddressRepository repository.UserAddressRepository, addressRepository repository.AddressRepository, transaction repository.Transaction) userAddressUsecaseImpl {
	return userAddressUsecaseImpl{
		userRepository:        userRepository,
		userAddressRepository: userAddressRepository,
		addressRepository:     addressRepository,
		transaction:           transaction,
	}
}

func (u *userAddressUsecaseImpl) UpdateUserAddress(ctx context.Context, accountId int64, updateUserAddressRequest dto.UpdateUserAddressRequest) error {
	newUserAddress := dto.ConvertUpdateRequestToUserAddress(updateUserAddressRequest)

	addressIdStr := ctx.Value(appconstant.AddressIdKey).(string)
	addressId, err := strconv.Atoi(addressIdStr)
	if err != nil {
		return apperror.AddressIdInvalidError()
	}

	if addressId < 1 {
		return apperror.AddressIdInvalidError()
	}

	newUserAddress.Id = int64(addressId)

	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if user == nil {
		return apperror.UserNotFoundError()
	}

	oldUserAddress, err := u.userAddressRepository.GetOneUserAddressByAddressId(ctx, int64(addressId))
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if oldUserAddress == nil {
		return apperror.UserAddressNotFoundError()
	}

	if user.Id != oldUserAddress.UserId {
		return apperror.ForbiddenAction()
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	userAddressRepo := tx.UserAddressRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if newUserAddress.IsMain && !oldUserAddress.IsMain {
		err = userAddressRepo.SetAllIsMainFalse(ctx, user.Id)
		if err != nil {
			return apperror.InternalServerError(err)
		}
	}

	err = userAddressRepo.UpdateOneUserAddress(ctx, newUserAddress)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *userAddressUsecaseImpl) AddUserAddress(ctx context.Context, accountId int64, addUserAddressRequest dto.AddUserAddressRequest) error {
	userAddress := dto.ConvertAddRequestToUserAddress(addUserAddressRequest)

	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if user == nil {
		return apperror.AccountNotFoundError()
	}
	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	userAddress.UserId = user.Id

	userAddressRepo := tx.UserAddressRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if userAddress.IsMain {
		err = userAddressRepo.SetAllIsMainFalse(ctx, userAddress.UserId)
		if err != nil {
			return apperror.InternalServerError(err)
		}
	}

	err = userAddressRepo.PostOneUserAddress(ctx, userAddress)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *userAddressUsecaseImpl) AddUserAddressAutofill(ctx context.Context, accountId int64, addUserAddressAutofillRequest dto.AddUserAddressAutofillRequest) error {
	userAddress := dto.ConvertAddRequestAutofillToUserAddress(addUserAddressAutofillRequest)

	provinceId, err := u.addressRepository.FindOneProvinceByName(ctx, userAddress.ProvinceName)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if provinceId == nil {
		return apperror.LocationNotFoundError()
	}

	userAddress.CityName = strings.ToLower(userAddress.CityName)
	if strings.Contains(userAddress.CityName, "kota") || strings.Contains(userAddress.CityName, "kabupaten") {
		userAddress.CityName = strings.ReplaceAll(userAddress.CityName, "kota", "")
		userAddress.CityName = strings.ReplaceAll(userAddress.CityName, "kabupaten", "")
	}
	cityId, err := u.addressRepository.FindOneCityByName(ctx, "%"+userAddress.CityName)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if cityId == nil {
		return apperror.LocationNotFoundError()
	}

	userAddress.DistrictName = strings.ToLower(userAddress.DistrictName)
	if strings.Contains(userAddress.DistrictName, "kecamatan") {
		userAddress.DistrictName = strings.ReplaceAll(userAddress.DistrictName, "kecamatan", "")
	}
	userAddress.DistrictName = strings.ReplaceAll(userAddress.DistrictName, " ", "%")
	districtId, err := u.addressRepository.FindOneDistrictByName(ctx, "%"+userAddress.DistrictName+"%")
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if districtId == nil {
		return apperror.LocationNotFoundError()
	}

	userAddress.SubdistrictName = strings.ToLower(userAddress.SubdistrictName)
	if strings.Contains(userAddress.SubdistrictName, "kelurahan") || strings.Contains(userAddress.SubdistrictName, "desa") {
		userAddress.SubdistrictName = strings.ReplaceAll(userAddress.SubdistrictName, "kelurahan", "")
		userAddress.SubdistrictName = strings.ReplaceAll(userAddress.SubdistrictName, "desa", "")
	}
	userAddress.SubdistrictName = strings.ReplaceAll(userAddress.SubdistrictName, " ", "%")
	subdistrictId, err := u.addressRepository.FindOneSubdistrictByName(ctx, userAddress.SubdistrictName)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if subdistrictId == nil {
		return apperror.LocationNotFoundError()
	}

	userAddress.Province.Id = *provinceId
	userAddress.City.Id = *cityId
	userAddress.District.Id = *districtId
	userAddress.Subdistrict.Id = *subdistrictId

	return u.AddUserAddress(ctx, accountId, dto.ConvertUserAddressToAddUserAddressRequest(userAddress))
}

func (u *userAddressUsecaseImpl) GetAllUserAddress(ctx context.Context, accountId int64) (*dto.AllUserAddressResponse, error) {
	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	userAddress, err := u.userAddressRepository.FindAllByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllUserAddressResponse(userAddress), nil
}

func (u *userAddressUsecaseImpl) DeleteUserAddress(ctx context.Context, userAddressId int64, accountId int64) error {
	userId, err := u.userAddressRepository.FindOneUserAddressById(ctx, userAddressId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if userId == nil {
		return apperror.UserAddressNotFoundError()
	}

	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if user == nil {
		return apperror.UserNotFoundError()
	}
	if user.Id != *userId {
		return apperror.ForbiddenAction()
	}

	err = u.userAddressRepository.DeleteOneUserAddress(ctx, userAddressId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
