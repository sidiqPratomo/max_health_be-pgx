package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/entity"
	"github.com/sidiqPratomo/max-health-backend/util"
)

func (u *authenticationUsecaseImpl) SendResetPasswordToken(ctx context.Context, sendEmailRequest dto.SendEmailRequest) error {
	acc, err := u.accountRepository.FindAccountByEmail(ctx, sendEmailRequest.Email)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if acc == nil {
		return apperror.NewAppError(http.StatusUnauthorized, errors.New("account not registered"), "account not registered")
	}

	if acc.VerifiedAt == nil {
		return apperror.AccountNotVerifiedError()
	}

	u.emailHelper.AddRequest([]string{sendEmailRequest.Email}, appconstant.ResetPasswordEmailSubject)

	resetPasswordToken, err := u.jwtHelper.CreateAndSign(util.JwtCustomClaims{
		UserId:        acc.Id,
		Email:         sendEmailRequest.Email,
		TokenDuration: 60,
	}, u.jwtHelper.Config.ResetPasswordSecret)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	resetPasswordUrl := fmt.Sprintf("http://%s%s/reset-password/verification/%s", os.Getenv("HOST"), os.Getenv("FE_PORT"), *resetPasswordToken)

	resetPasswordCode := util.GenerateCode(6)

	err = u.emailHelper.CreateBody(appconstant.ResetPasswordEmailTemplate, struct {
		Name string
		Url  string
		Code string
	}{
		Name: acc.Name,
		Url:  resetPasswordUrl,
		Code: resetPasswordCode,
	})
	if err != nil {
		return apperror.InternalServerError(err)
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	resetPasswordTokenRepo := tx.ResetPasswordTokenRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	err = resetPasswordTokenRepo.InvalidateTokens(ctx, acc.Id)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = resetPasswordTokenRepo.PostOneToken(ctx, acc.Id, resetPasswordCode)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	err = u.emailHelper.SendEmail()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}

func (u *authenticationUsecaseImpl) ResetPassword(ctx context.Context, resetPasswordTokenVerificationRequest dto.ResetPasswordVerificationRequest) error {
	resetPasswordToken := dto.ResetPasswordVerificationRequestToResetPasswordToken(resetPasswordTokenVerificationRequest)

	resetPasswordToken, err := u.resetPasswordTokenRepository.FindOneByAccountIdAndToken(ctx, resetPasswordToken)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	if resetPasswordToken == nil {
		return apperror.InvalidCodeError()
	}
	if time.Now().After(resetPasswordToken.ExpiredAt.Add(-7 * time.Hour)) {
		return apperror.ExpiredCodeError()
	}

	oldPassword, err := u.accountRepository.FindOnePasswordById(ctx, resetPasswordToken.AccountId)
	if err != nil {
		return apperror.InternalServerError(err)
	}
	if oldPassword == nil {
		return apperror.AccountNotFoundError()
	}

	isPasswordReused, err := u.hashHelper.CheckPassword(resetPasswordTokenVerificationRequest.Password, []byte(*oldPassword))
	if isPasswordReused {
		return apperror.OldPasswordReusedError()
	}

	newPassword, err := u.hashHelper.HashPassword(resetPasswordTokenVerificationRequest.Password)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	tx, err := u.transaction.BeginTx()
	if err != nil {
		return apperror.InternalServerError(err)
	}

	accountRepo := tx.AccountRepository()
	resetPasswordTokenRepo := tx.ResetPasswordTokenRepository()

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err = resetPasswordTokenRepo.UpdateExpiredAtOne(ctx, resetPasswordToken); err != nil {
		return apperror.InternalServerError(err)
	}

	if err = accountRepo.UpdatePasswordOne(ctx, &entity.Account{Id: resetPasswordToken.AccountId, Password: newPassword}); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
