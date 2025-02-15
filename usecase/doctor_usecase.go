package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type DoctorUsecase interface {
	UpdateData(
		ctx context.Context,
		account entity.DetailedDoctor,
		file multipart.File,
		fileHeader *multipart.FileHeader) error
	GetAllDoctors(
		c context.Context,
		Sort string,
		SortBy string,
		Limit string,
		SpecializationId string,
		Page string) (*dto.GetAllDoctorResponse, error)
	GetAllDoctorSpecialization(ctx context.Context) ([]dto.DoctorSpecialization, error)
	GetProfile(ctx context.Context, accountId int64) (*dto.DoctorProfileResponse, error)
	GetProfileForPublic(ctx context.Context, doctorId int64) (*dto.DoctorProfileResponse, error)
	UpdateDoctorStatus(ctx context.Context, doctorAccountId int64, isOnline bool) error
	GetDoctorIsOnline(ctx context.Context, doctorAccountId int64) (*dto.GetDoctorStatusResponse, error)
}

type doctorUsecaseImpl struct {
	accountRepository              repository.AccountRepository
	doctorRepository               repository.DoctorRepository
	doctorSpecializationRepository repository.DoctorSpecializationRepository
	transaction                    repository.Transaction
	hashHelper                     util.HashHelperIntf
}

func NewDoctorUsecaseImpl(accountRepository repository.AccountRepository, doctorRepository repository.DoctorRepository, doctorSpecializationRepository repository.DoctorSpecializationRepository, transaction repository.Transaction, hashHelper util.HashHelperIntf) doctorUsecaseImpl {
	return doctorUsecaseImpl{
		accountRepository:              accountRepository,
		doctorRepository:               doctorRepository,
		doctorSpecializationRepository: doctorSpecializationRepository,
		transaction:                    transaction,
		hashHelper:                     hashHelper,
	}
}

func (u *doctorUsecaseImpl) UpdateData(ctx context.Context, doctor entity.DetailedDoctor, file multipart.File, fileHeader *multipart.FileHeader) error {
	dbAccount, err := u.accountRepository.FindOneById(ctx, doctor.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if dbAccount == nil {
		return apperror.EmailNotFoundError()
	}

	doctor.Name = strings.Trim(doctor.Name, " ")

	isNameValid := util.RegexValidate(doctor.Name, appconstant.NameRegexPattern)
	if !isNameValid {
		return apperror.InvalidNameError(errors.New(("invalid name")))
	}

	if doctor.Password != "" {
		if !util.ValidatePassword(doctor.Password) {
			return apperror.NewAppError(http.StatusBadRequest, errors.New("invalid password"), "invalid password")
		}
	}

	if file != nil {
		filePath, _, err := util.ValidateFile(*fileHeader, appconstant.ProfilePicturesUrl, []string{"png", "jpg", "jpeg"}, 2000000)
		if err != nil {
			return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
		}

		imageUrl, err := util.UploadToCloudinary(file, *filePath)
		if err != nil {
			return apperror.InternalServerError(err)
		}
		doctor.ProfilePicture = imageUrl
		if dbAccount.ProfilePicture != "https://res.cloudinary.com/dpdu3tidt/image/upload/v1713774687/profile_pictures/xdv5xzkz1yr0qwgkc6yk.avif" {
			util.DeleteInCloudinary(dbAccount.ProfilePicture)
		}
	} else {
		doctor.ProfilePicture = dbAccount.ProfilePicture
	}

	if doctor.Password != "" {
		password, err := u.hashHelper.HashPassword(doctor.Password)
		if err != nil {
			return apperror.InternalServerError(err)
		}
		doctor.Password = password
	} else {
		doctor.Password = dbAccount.Password
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	doctorRepo := tx.DoctorRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	var account = entity.Account{Id: doctor.Id, Name: doctor.Name, Password: doctor.Password, ProfilePicture: doctor.ProfilePicture}
	err = accountRepo.UpdateDataOne(ctx, &account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = doctorRepo.UpdateDataOne(ctx, &doctor)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *doctorUsecaseImpl) GetAllDoctors(
	c context.Context,
	Sort string,
	SortBy string,
	Limit string,
	SpecializationId string,
	Page string) (*dto.GetAllDoctorResponse, error) {

	pageInt, err := strconv.Atoi(Page)
	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(Limit)
	if err != nil {
		limitInt = 12
	}

	offset := (pageInt - 1) * limitInt
	if offset < 0 {
		offset = 0
	}

	SortList := strings.Split(Sort, ",")
	SortByList := strings.Split(SortBy, ",")

	if Sort == "" && SortBy != "" {
		return nil, apperror.BadRequestError(errors.New("sort parameter is required when sortBy is provided"))
	}

	if SortBy == "" && Sort != "" {
		return nil, apperror.BadRequestError(errors.New("sortBy parameter is required when sort is provided"))
	}

	if Sort != "" {
		for _, sortDirection := range SortList {
			if sortDirection != "asc" && sortDirection != "desc" {
				return nil, apperror.BadRequestError(errors.New("invalid sort direction"))
			}
		}
	}

	if SortBy != "" {
		for _, sortByName := range SortByList {
			if sortByName != "fee_per_patient" && sortByName != "experience" {
				return nil, apperror.BadRequestError(errors.New("invalid sortBy Name"))
			}
		}
	}

	if len(SortByList) != len(SortList) {
		return nil, apperror.BadRequestError(errors.New("invalid length Sort or SortBy"))
	}

	doctors, pageInfo, err := u.doctorRepository.GetAllDoctor(c, SortList, SortByList, Limit, offset, SpecializationId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	getAllDoctor := []dto.DoctorDto{}

	for _, doctor := range doctors {
		account, err := u.accountRepository.FindOneById(c, doctor.AccountId)
		if err != nil {
			return nil, err
		}
		specialist, err := u.doctorRepository.FindSpecializationById(c, doctor.SpecializationId)
		if err != nil {
			return nil, err
		}

		doctorDto := dto.DoctorDto{
			DoctorId:           doctor.Id,
			AccountId:          doctor.AccountId,
			FeePerPatient:      doctor.FeePerPatient,
			IsOnline:           doctor.IsOnline,
			ProfilePicture:     account.ProfilePicture,
			Experience:         doctor.Experience,
			Name:               account.Name,
			SpecializationName: *specialist,
		}

		getAllDoctor = append(getAllDoctor, doctorDto)
	}
	doctorResponse := dto.GetAllDoctorResponse{
		PageInfo: *pageInfo,
		Doctors:  getAllDoctor,
	}
	return &doctorResponse, nil
}

func (u *doctorUsecaseImpl) GetAllDoctorSpecialization(ctx context.Context) ([]dto.DoctorSpecialization, error) {
	doctorSpecializationList, err := u.doctorSpecializationRepository.GetAllDoctorSpecialization(ctx)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	doctorSpecializationListDTO := dto.ConvertToDoctorSpecializationDTOList(doctorSpecializationList)

	return doctorSpecializationListDTO, nil
}

func (u *doctorUsecaseImpl) GetProfile(ctx context.Context, accountId int64) (*dto.DoctorProfileResponse, error) {
	account, err := u.accountRepository.FindOneById(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if account == nil {
		return nil, apperror.AccountNotFoundError()
	}

	doctor, err := u.doctorRepository.FindDoctorByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if doctor == nil {
		return nil, apperror.DoctorNotFoundError()
	}

	res := &dto.DoctorProfileResponse{
		Email:          account.Email,
		Name:           account.Name,
		ProfilePicture: account.ProfilePicture,
	}

	res.SpecializationName = doctor.SpecializationName
	res.SpecializationId = doctor.SpecializationId
	res.FeePerPatient = doctor.FeePerPatient
	res.Experience = doctor.Experience

	return res, nil
}

func (u *doctorUsecaseImpl) GetProfileForPublic(ctx context.Context, doctorId int64) (*dto.DoctorProfileResponse, error) {
	doctor, err := u.doctorRepository.FindDoctorByDoctorId(ctx, doctorId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if doctor == nil {
		return nil, apperror.DoctorNotFoundError()
	}

	res := &dto.DoctorProfileResponse{
		Email:          doctor.Email,
		Name:           doctor.Name,
		ProfilePicture: doctor.ProfilePicture,
	}

	res.SpecializationName = doctor.SpecializationName
	res.SpecializationId = doctor.SpecializationId
	res.FeePerPatient = doctor.FeePerPatient
	res.Experience = doctor.Experience

	return res, nil
}

func (u *doctorUsecaseImpl) UpdateDoctorStatus(ctx context.Context, doctorAccountId int64, isOnline bool) error {
	doctor, err := u.doctorRepository.FindDoctorByAccountId(ctx, doctorAccountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if doctor == nil {
		return apperror.ForbiddenAction()
	}

	err = u.doctorRepository.UpdateDoctorStatus(ctx, doctorAccountId, isOnline)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *doctorUsecaseImpl) GetDoctorIsOnline(ctx context.Context, doctorAccountId int64) (*dto.GetDoctorStatusResponse, error) {
	doctor, err := u.doctorRepository.FindDoctorByAccountId(ctx, doctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if doctor == nil {
		return nil, apperror.DoctorNotFoundError()
	}

	isOnline, err := u.doctorRepository.GetDoctorIsOnline(ctx, doctorAccountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &dto.GetDoctorStatusResponse{IsOnline: *isOnline}, nil
}
