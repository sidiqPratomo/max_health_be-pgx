package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type PartnerUsecase interface {
	AddOnePartner(ctx context.Context, registerRequest dto.RegisterRequest, file multipart.File, fileHeader multipart.FileHeader) error
	GetAllPartners(ctx context.Context) (*dto.AllPartnersResponse, error)
	UpdateOnePartner(ctx context.Context, updateAccountRequest dto.UpdateAccountRequest, pharmacyManagerId int64, file multipart.File, fileHeader *multipart.FileHeader) error
	DeleteOnePartner(ctx context.Context, pharmacyManagerId int64) error
	SendCredentials(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error
}

type partnerUsecaseImpl struct {
	accountRepository         repository.AccountRepository
	pharmacyManagerRepository repository.PharmacyManagerRepository
	transaction               repository.Transaction
	hashHelper                util.HashHelperIntf
	emailHelper               util.EmailHelper
}

type PartnerUsecaseImplOpts struct {
	AccountRepository         repository.AccountRepository
	PharmacyManagerRepository repository.PharmacyManagerRepository
	Transaction               repository.Transaction
	HashHelper                util.HashHelperIntf
	EmailHelper               util.EmailHelper
}

func NewPartnerUsecaseImpl(opts PartnerUsecaseImplOpts) partnerUsecaseImpl {
	return partnerUsecaseImpl{
		accountRepository:         opts.AccountRepository,
		pharmacyManagerRepository: opts.PharmacyManagerRepository,
		transaction:               opts.Transaction,
		hashHelper:                opts.HashHelper,
		emailHelper:               opts.EmailHelper,
	}
}

func (u *partnerUsecaseImpl) AddOnePartner(ctx context.Context, registerRequest dto.RegisterRequest, file multipart.File, fileHeader multipart.FileHeader) error {
	account := dto.RegisterRequestToAccount(registerRequest)

	account.Name = strings.Trim(account.Name, " ")

	if isNameValid := util.RegexValidate(account.Name, appconstant.NameRegexPattern); !isNameValid {
		return apperror.InvalidNameError(errors.New("invalid name"))
	}

	filePath, _, err := util.ValidateFile(fileHeader, appconstant.ProfilePicturesUrl, []string{"jpg", "jpeg", "png"}, 2000000)
	if err != nil {
		return apperror.BadRequestError(err)
	}

	acc, err := u.accountRepository.FindAccountByEmail(ctx, account.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if acc != nil {
		return apperror.EmailTakenError()
	}

	imageUrl, err := util.UploadToCloudinary(file, *filePath)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	password, err := u.hashHelper.HashPassword(util.GenerateCode(8))
	if err != nil {
		return apperror.InternalServerError(err)
	}

	account.Password = password
	account.RoleId = appconstant.PharmacyManagerId
	account.ProfilePicture = imageUrl

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	pharmacyManagerRepo := tx.PharmacyManagerRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	accountId, err := accountRepo.PostOneVerifiedAccount(ctx, account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if err = pharmacyManagerRepo.PostOne(ctx, *accountId); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *partnerUsecaseImpl) GetAllPartners(ctx context.Context) (*dto.AllPartnersResponse, error) {
	pharmacyManagers, err := u.pharmacyManagerRepository.FindAll((ctx))
	if err != nil {
		return nil, err
	}

	return dto.ConvertToAllPartnersResponse(pharmacyManagers), err
}

func (u *partnerUsecaseImpl) UpdateOnePartner(ctx context.Context, updateAccountRequest dto.UpdateAccountRequest, pharmacyManagerId int64, file multipart.File, fileHeader *multipart.FileHeader) error {
	accountRequest := dto.UpdateAccountRequestToAccount(updateAccountRequest)

	accountRequest.Name = strings.Trim(accountRequest.Name, " ")

	isNameValid := util.RegexValidate(accountRequest.Name, appconstant.NameRegexPattern)
	if !isNameValid {
		return apperror.InvalidNameError(errors.New((appconstant.MsgInvalidName)))
	}

	pharmacyManager, err := u.pharmacyManagerRepository.FindOneById(ctx, pharmacyManagerId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if pharmacyManager == nil {
		return apperror.PartnerNotFoundError()
	}

	account, err := u.accountRepository.FindOneById(ctx, pharmacyManager.Account.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if account == nil {
		return apperror.AccountNotFoundError()
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
		accountRequest.ProfilePicture = imageUrl
		if account.ProfilePicture != "https://res.cloudinary.com/dpdu3tidt/image/upload/v1713774687/profile_pictures/xdv5xzkz1yr0qwgkc6yk.avif" {
			util.DeleteInCloudinary(account.ProfilePicture)
		}
	} else {
		accountRequest.ProfilePicture = account.ProfilePicture
	}

	accountRequest.Id = pharmacyManager.Account.Id
	if err = u.accountRepository.UpdateNameAndProfilePictureOne(ctx, &accountRequest); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *partnerUsecaseImpl) DeleteOnePartner(ctx context.Context, pharmacyManagerId int64) error {
	pharmacyManager, err := u.pharmacyManagerRepository.FindOneById(ctx, pharmacyManagerId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if pharmacyManager == nil {
		return apperror.PartnerNotFoundError()
	}

	if strings.Split(pharmacyManager.Account.ProfilePicture, "/")[1] == "res.cloudinary.com" {
		util.DeleteInCloudinary(pharmacyManager.Account.ProfilePicture)
	}

	tx, err := u.transaction.BeginTx(ctx)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	pharmacyManagerRepo := tx.PharmacyManagerRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err = pharmacyManagerRepo.DeleteOneById(ctx, pharmacyManagerId); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = accountRepo.DeleteOneById(ctx, pharmacyManager.Account.Id); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *partnerUsecaseImpl) SendCredentials(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error {
	account, err := u.accountRepository.FindAccountByEmail(ctx, sendEmailRequest.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if account == nil {
		return apperror.NewAppError(http.StatusUnauthorized, errors.New(appconstant.MsgAccountNotRegistered), appconstant.MsgAccountNotRegistered)
	}

	rawPassword := util.GenerateCode(8)
	password, err := u.hashHelper.HashPassword(rawPassword)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	account.Password = password

	u.emailHelper.AddRequest([]string{sendEmailRequest.Email}, appconstant.CredentialsEmailSubject)

	if err = u.emailHelper.CreateBody(appconstant.CredentialsEmailTemplate, struct {
		Name        string
		Email       string
		Credentials string
	}{
		Name:        account.Name,
		Email:       account.Email,
		Credentials: rawPassword,
	}); err != nil {
		return apperror.InternalServerError(err)
	}

	err = u.accountRepository.UpdatePasswordOne(ctx, account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = u.emailHelper.SendEmail()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
