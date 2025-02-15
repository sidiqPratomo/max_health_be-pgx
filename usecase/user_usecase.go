package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
)

type UserUsecase interface {
	GetProfile(ctx context.Context, accountId int64) (*dto.UserProfileResponse, error)
	UpdateData(ctx context.Context, user entity.DetailedUser, file multipart.File, fileHeader *multipart.FileHeader) error
}

type userUsecaseImpl struct {
	accountRepository     repository.AccountRepository
	transaction           repository.Transaction
	userRepository        repository.UserRepository
	userAddressRepository repository.UserAddressRepository
	hashHelper            util.HashHelperIntf
}

func NewUserUsecaseImpl(accountRepository repository.AccountRepository, transaction repository.Transaction, userRepository repository.UserRepository, userAddressRepository repository.UserAddressRepository, hashHelper util.HashHelperIntf) userUsecaseImpl {
	return userUsecaseImpl{
		accountRepository:     accountRepository,
		transaction:           transaction,
		userRepository:        userRepository,
		userAddressRepository: userAddressRepository,
		hashHelper:            hashHelper,
	}
}

func (u *userUsecaseImpl) UpdateData(ctx context.Context, user entity.DetailedUser, file multipart.File, fileHeader *multipart.FileHeader) error {
	dbAccount, err := u.accountRepository.FindOneById(ctx, user.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if dbAccount == nil {
		return apperror.EmailNotFoundError()
	}

	date, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		return apperror.BadRequestError(errors.New("invalid date"))
	}
	if date.After(time.Now()) {
		return apperror.BadRequestError(errors.New("date of birth can't be later than current date"))
	}
	user.Name = strings.Trim(user.Name, " ")

	isNameValid := util.RegexValidate(user.Name, appconstant.NameRegexPattern)
	if !isNameValid {
		return apperror.InvalidNameError(errors.New("invalid name"))
	}

	if user.Password != "" {
		if !util.ValidatePassword(user.Password) {
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
		user.ProfilePicture = imageUrl
		if dbAccount.ProfilePicture != "https://res.cloudinary.com/dpdu3tidt/image/upload/v1713774687/profile_pictures/xdv5xzkz1yr0qwgkc6yk.avif" {
			util.DeleteInCloudinary(dbAccount.ProfilePicture)
		}
	} else {
		user.ProfilePicture = dbAccount.ProfilePicture
	}

	if user.Password != "" {
		password, err := u.hashHelper.HashPassword(user.Password)
		if err != nil {
			return apperror.InternalServerError(err)
		}
		user.Password = password
	} else {
		user.Password = dbAccount.Password
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	userRepo := tx.UserRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	var account = entity.Account{Id: user.Id, Name: user.Name, Password: user.Password, ProfilePicture: user.ProfilePicture}
	err = accountRepo.UpdateDataOne(ctx, &account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = userRepo.UpdateDataOne(ctx, &user)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *userUsecaseImpl) GetProfile(ctx context.Context, accountId int64) (*dto.UserProfileResponse, error) {
	account, err := u.accountRepository.FindOneById(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if account == nil {
		return nil, apperror.AccountNotFoundError()
	}

	user, err := u.userRepository.FindUserByAccountId(ctx, accountId)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if user == nil {
		return nil, apperror.UserNotFoundError()
	}

	res := &dto.UserProfileResponse{
		Email:          account.Email,
		Name:           account.Name,
		ProfilePicture: account.ProfilePicture,
	}
	if user.GenderName != nil {
		res.Gender = *user.GenderName
		res.GenderId = *user.GenderId
	}
	if user.DateOfBirth != nil {
		res.DateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	return res, nil
}
