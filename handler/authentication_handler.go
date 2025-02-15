package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthenticationHandler struct {
	authenticationUsecase usecase.AuthenticationUsecase
}

func NewAuthenticationHandler(authenticationUsecase usecase.AuthenticationUsecase) AuthenticationHandler {
	return AuthenticationHandler{
		authenticationUsecase: authenticationUsecase,
	}
}

func (h *AuthenticationHandler) Login(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	var loginRequest dto.LoginRequest
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	account := dto.LoginRequestToAccount(loginRequest)

	tokens, err := h.authenticationUsecase.Login(ctx.Request.Context(), account)
	if err != nil {
		ctx.Error(err)
		return
	}

	resTokens := dto.ConvertTokensToResponse(*tokens)
	util.ResponseOK(ctx, resTokens)
}

func (h *AuthenticationHandler) RegisterUser(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var registerRequest dto.RegisterRequest

	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.authenticationUsecase.RegisterUser(ctx.Request.Context(), registerRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *AuthenticationHandler) RegisterDoctor(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var registerDoctorRequest dto.RegisterDoctorRequest
	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("data is empty"), "data is empty"))
		return
	}

	err := json.Unmarshal([]byte(data), &registerDoctorRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(registerDoctorRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file == nil {
			ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("file not attached"), "file not attached"))
			return
		}
		ctx.Error(err)
		return
	}

	var registerRequest dto.RegisterRequest
	registerRequest.Email = registerDoctorRequest.Email
	registerRequest.Name = registerDoctorRequest.Name

	err = h.authenticationUsecase.RegisterDoctor(ctx.Request.Context(), registerRequest, registerDoctorRequest.SpecializationId, file, *fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *AuthenticationHandler) SendVerificationEmail(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var sendEmailRequest dto.SendEmailRequest

	err := ctx.ShouldBindJSON(&sendEmailRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.authenticationUsecase.SendVerificationEmail(ctx.Request.Context(), sendEmailRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *AuthenticationHandler) GetNewAccessToken(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var accessTokenRequest dto.AccessTokenRequest

	err := ctx.ShouldBindJSON(&accessTokenRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	accessToken, err := h.authenticationUsecase.GetNewAccessToken(ctx.Request.Context(), accessTokenRequest.RefreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	accessTokenRes := dto.ConvertAccessTokenToResponse(accessToken)
	util.ResponseOK(ctx, accessTokenRes)
}

func (h *AuthenticationHandler) VerifyOneAccount(ctx *gin.Context) {
	var verificationPasswordRequest dto.VerificationPasswordRequest

	if err := ctx.ShouldBindJSON(&verificationPasswordRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.authenticationUsecase.VerifyOneAccount(ctx.Request.Context(), verificationPasswordRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *AuthenticationHandler) SendResetPasswordToken(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var sendEmailRequest dto.SendEmailRequest

	if err := ctx.ShouldBindJSON(&sendEmailRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.authenticationUsecase.SendResetPasswordToken(ctx.Request.Context(), sendEmailRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *AuthenticationHandler) ResetPasswordOneAccount(ctx *gin.Context) {
	var resetPasswordTokenRequest dto.ResetPasswordVerificationRequest

	if err := ctx.ShouldBindJSON(&resetPasswordTokenRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.authenticationUsecase.ResetPassword(ctx.Request.Context(), resetPasswordTokenRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}
