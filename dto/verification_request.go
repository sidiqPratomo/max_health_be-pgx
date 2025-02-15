package dto

import "github.com/sidiqPratomo/max-health-backend/entity"

type VerificationPasswordRequest struct {
	AccountId        int64  `json:"account_id" binding:"required"`
	Password         string `json:"password" binding:"required,ValidPassword"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

func VerificationPasswordRequestToVerificationCode(verificationPasswordRequest VerificationPasswordRequest) *entity.VerificationCode {
	return &entity.VerificationCode{
		AccountId: verificationPasswordRequest.AccountId,
		Code:      verificationPasswordRequest.VerificationCode,
	}
}
