package dto

import (
	"github.com/sidiqPratomo/max-health-backend/entity"
)

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
	Name  string `json:"name" binding:"required" validate:"required"`
}

type RegisterDoctorRequest struct {
	Email            string `json:"email" binding:"required,email" validate:"required,email"`
	Name             string `json:"name" binding:"required" validate:"required"`
	SpecializationId int64  `json:"specialization_id" binding:"required" validate:"required,gte=1,lte=17"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type SendEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordVerificationRequest struct {
	AccountId int64  `json:"account_id" binding:"required"`
	Password  string `json:"password" binding:"required,ValidPassword"`
	Code      string `json:"code" binding:"required"`
}

type UpdateDataRequest struct {
	Name              string `json:"name" validate:"required"`
	Password          string `json:"password"`
	FeePerPatient     int    `json:"fee_per_patient" validate:"number,gte=0"`
	GenderId          int64  `json:"gender_id" validate:"omitempty,number,gte=1"`
	DateOfBirth       string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
	YearsOfExperience int    `json:"years_of_experience" validate:"omitempty,number,gte=0"`
}

type UpdateAccountRequest struct {
	Name string `json:"name" validate:"required"`
}

func RegisterRequestToAccount(RegisterRquest RegisterRequest) entity.Account {
	return entity.Account{
		Email: RegisterRquest.Email,
		Name:  RegisterRquest.Name,
	}
}

func LoginRequestToAccount(LoginRequest LoginRequest) entity.Account {
	return entity.Account{
		Email:    LoginRequest.Email,
		Password: LoginRequest.Password,
	}
}

func ResetPasswordVerificationRequestToResetPasswordToken(resetPasswordVerificationRequest ResetPasswordVerificationRequest) *entity.ResetPasswordToken {
	return &entity.ResetPasswordToken{
		AccountId: resetPasswordVerificationRequest.AccountId,
		Token:     resetPasswordVerificationRequest.Code,
	}
}

func UpdateAccountRequestToAccount(updateAccountRequest UpdateAccountRequest) entity.Account {
	return entity.Account{
		Name: updateAccountRequest.Name,
	}
}
