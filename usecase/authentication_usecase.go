package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/repository"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/google/uuid"
)

type AuthenticationUsecase interface {
	RegisterDoctor(ctx context.Context, registerRequest dto.RegisterRequest, specializationId int64, file multipart.File, fileHeader multipart.FileHeader) error
	RegisterUser(ctx context.Context, registerRequest dto.RegisterRequest) error
	Login(ctx context.Context, account entity.Account) (*entity.Tokens, error)
	SendVerificationEmail(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error
	GetNewAccessToken(ctx context.Context, refreshToken string) (string, error)
	VerifyOneAccount(ctx context.Context, verificationPasswordRequest dto.VerificationPasswordRequest) error
	SendResetPasswordToken(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error
	ResetPassword(ctx context.Context, resetPasswordTokenVerificationRequest dto.ResetPasswordVerificationRequest) error
}

type authenticationUsecaseImpl struct {
	accountRepository            repository.AccountRepository
	userRepository               repository.UserRepository
	doctorRepository             repository.DoctorRepository
	verificationCodeRepository   repository.VerificationCodeRepository
	refreshTokenRepository       repository.RefreshTokenRepository
	resetPasswordTokenRepository repository.ResetPasswordTokenRepository
	transaction                  repository.Transaction
	hashHelper                   util.HashHelperIntf
	jwtHelper                    util.JwtAuthentication
	emailHelper                  util.EmailHelper
}

type AuthenticationUsecaseImplOpts struct {
	DrugRepository               repository.DrugRepository
	PharmacyDrugRepository       repository.PharmacyDrugRepository
	CartRepository               repository.CartRepository
	AccountRepository            repository.AccountRepository
	UserRepository               repository.UserRepository
	DoctorRepository             repository.DoctorRepository
	VerificationCodeRepository   repository.VerificationCodeRepository
	RefreshTokenRepositoy        repository.RefreshTokenRepository
	ResetPasswordTokenRepository repository.ResetPasswordTokenRepository
	Transaction                  repository.Transaction
	HashHelper                   util.HashHelperIntf
	JwtHelper                    util.JwtAuthentication
	EmailHelper                  util.EmailHelper
}

func NewAuthenticationUsecaseImpl(opts AuthenticationUsecaseImplOpts) authenticationUsecaseImpl {
	return authenticationUsecaseImpl{
		accountRepository:            opts.AccountRepository,
		userRepository:               opts.UserRepository,
		doctorRepository:             opts.DoctorRepository,
		verificationCodeRepository:   opts.VerificationCodeRepository,
		refreshTokenRepository:       opts.RefreshTokenRepositoy,
		resetPasswordTokenRepository: opts.ResetPasswordTokenRepository,
		transaction:                  opts.Transaction,
		hashHelper:                   opts.HashHelper,
		jwtHelper:                    opts.JwtHelper,
		emailHelper:                  opts.EmailHelper,
	}
}

func (u *authenticationUsecaseImpl) RegisterDoctor(ctx context.Context, registerRequest dto.RegisterRequest, specializationId int64, file multipart.File, fileHeader multipart.FileHeader) error {
	account := dto.RegisterRequestToAccount(registerRequest)

	acc, err := u.accountRepository.FindAccountByEmail(ctx, account.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if acc != nil {
		return apperror.NewAppError(http.StatusConflict, errors.New("email has been taken"), "email has been taken")
	}

	account.Name = strings.Trim(account.Name, " ")

	isNameValid := util.RegexValidate(account.Name, appconstant.NameRegexPattern)
	if !isNameValid {
		return apperror.InvalidNameError(errors.New(("invalid name")))
	}

	specialization, err := u.doctorRepository.FindSpecializationById(ctx, specializationId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if specialization == nil {
		return apperror.NewAppError(http.StatusBadRequest, errors.New("invalid specialization id"), "invalid specialization id")
	}

	filePath, _, err := util.ValidateFile(fileHeader, appconstant.DoctorCertificatesUrl, []string{"pdf"}, 2000000)
	if err != nil {
		return apperror.NewAppError(http.StatusBadRequest, err, err.Error())
	}

	imageUrl, err := util.UploadToCloudinary(file, *filePath)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	password, err := u.hashHelper.HashPassword(uuid.NewString())
	if err != nil {
		return apperror.InternalServerError(err)
	}

	account.Password = password
	account.RoleId = appconstant.DoctorId

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

	accountId, err := accountRepo.PostOneAccount(ctx, account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = doctorRepo.PostOneDoctor(ctx, *accountId, specializationId, imageUrl)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *authenticationUsecaseImpl) RegisterUser(ctx context.Context, registerRequest dto.RegisterRequest) error {
	account := dto.RegisterRequestToAccount(registerRequest)

	acc, err := u.accountRepository.FindAccountByEmail(ctx, account.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if acc != nil {
		return apperror.NewAppError(http.StatusConflict, errors.New("email has been taken"), "email has been taken")
	}

	isNameValid := util.RegexValidate(account.Name, appconstant.NameRegexPattern)
	if !isNameValid {
		return apperror.InvalidNameError(errors.New("invalid name"))
	}

	password, err := u.hashHelper.HashPassword(uuid.NewString())
	if err != nil {
		return apperror.InternalServerError(err)
	}

	account.Password = password
	account.RoleId = appconstant.UserId

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

	accountId, err := accountRepo.PostOneAccount(ctx, account)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = userRepo.PostOneUser(ctx, *accountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *authenticationUsecaseImpl) SendVerificationEmail(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error {
	acc, err := u.accountRepository.FindAccountByEmail(ctx, sendEmailRequest.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if acc == nil {
		return apperror.NewAppError(http.StatusUnauthorized, errors.New("account not registered"), "account not registered")
	}

	if acc.VerifiedAt != nil {
		return apperror.NewAppError(http.StatusForbidden, errors.New("account has been verified"), "account has been verified")
	}

	u.emailHelper.AddRequest([]string{sendEmailRequest.Email}, appconstant.VerificationEmailSubject)

	verificationToken, err := u.jwtHelper.CreateAndSign(util.JwtCustomClaims{
		UserId:        acc.Id,
		Email:         sendEmailRequest.Email,
		TokenDuration: 60,
	}, u.jwtHelper.Config.VerifSecret)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	verificationUrl := fmt.Sprintf("https://%s/auth/verification/%s", os.Getenv("EMAILHOST"), *verificationToken)

	verificationCode := util.GenerateCode(6)

	err = u.emailHelper.CreateBody(appconstant.VerificationEmailTemplate, struct {
		Name string
		Url  string
		Code string
	}{
		Name: acc.Name,
		Url:  verificationUrl,
		Code: verificationCode,
	})
	if err != nil {
		return apperror.InternalServerError(err)
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	verificationCodeRepo := tx.VerificationCodeRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	err = verificationCodeRepo.InvalidateCodes(ctx, acc.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = verificationCodeRepo.PostOneCode(ctx, acc.Id, verificationCode)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = u.emailHelper.SendEmail()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *authenticationUsecaseImpl) Login(ctx context.Context, account entity.Account) (*entity.Tokens, error) {
	tokens := entity.Tokens{}
	userCredential, err := u.accountRepository.FindAccountByEmail(ctx, account.Email)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if userCredential == nil {
		return nil, apperror.EmailNotFoundError()
	}

	if userCredential.VerifiedAt == nil {
		return nil, apperror.AccountNotVerifiedError()
	}

	isPassword, err := u.hashHelper.CheckPassword(account.Password, []byte(userCredential.Password))
	if !isPassword {
		return nil, apperror.WrongPasswordError(err)
	}

	customClaims := util.JwtCustomClaims{UserId: userCredential.Id, Email: userCredential.Email, Role: userCredential.RoleName, TokenDuration: 15}
	accessToken, err := u.jwtHelper.CreateAndSign(customClaims, u.jwtHelper.Config.AccessSecret)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	customClaims.TokenDuration = 24 * 60
	refreshToken, err := u.jwtHelper.CreateAndSign(customClaims, u.jwtHelper.Config.RefreshSecret)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	tokens.AccessToken = *accessToken
	tokens.RefreshToken = *refreshToken

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	refreshTokenRepo := tx.RefreshTokenRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	err = refreshTokenRepo.InvalidateCodes(ctx, userCredential.Id)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	err = refreshTokenRepo.PostOneCode(ctx, userCredential.Id, *refreshToken)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &tokens, nil
}

func (u *authenticationUsecaseImpl) GetNewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claim, err := u.jwtHelper.ParseAndVerify(refreshToken, u.jwtHelper.Config.RefreshSecret)
	if err != nil {
		return "", apperror.InternalServerError(err)
	}

	refreshToken, err = u.refreshTokenRepository.FindOneCode(ctx, refreshToken)
	if err != nil {
		return "", apperror.InternalServerError(err)
	}
	if refreshToken == "" {
		return "", apperror.RefreshTokenExpiredError()
	}

	claim.TokenDuration = 15
	accessToken, err := u.jwtHelper.CreateAndSign(*claim, u.jwtHelper.Config.AccessSecret)
	if err != nil {
		return "", apperror.InternalServerError(err)
	}

	return *accessToken, nil
}

func (u *authenticationUsecaseImpl) VerifyOneAccount(ctx context.Context, verificationPasswordRequest dto.VerificationPasswordRequest) error {
	verificationCode := dto.VerificationPasswordRequestToVerificationCode(verificationPasswordRequest)

	verificationCode, err := u.verificationCodeRepository.FindOneByAccountIdAndCode(ctx, verificationCode)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if verificationCode == nil {
		return apperror.InvalidCodeError()
	}
	if time.Now().After(verificationCode.ExpiredAt.Add(-7 * time.Hour)) {
		return apperror.ExpiredCodeError()
	}

	password, err := u.hashHelper.HashPassword(verificationPasswordRequest.Password)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	verificationRepo := tx.VerificationCodeRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err = verificationRepo.UpdateExpiredAtOne(ctx, verificationCode); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = accountRepo.UpdatePasswordOne(ctx, &entity.Account{Id: verificationCode.AccountId, Password: password}); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
