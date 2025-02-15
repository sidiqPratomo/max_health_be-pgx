package handler

import (
	"encoding/json"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PartnerHandler struct {
	partnerUsecase usecase.PartnerUsecase
}

func NewPartnerHandler(partnerUsecase usecase.PartnerUsecase) PartnerHandler {
	return PartnerHandler{
		partnerUsecase: partnerUsecase,
	}
}

func (h *PartnerHandler) AddPartner(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var registerRequest dto.RegisterRequest

	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.EmptyDataError())
		return
	}

	if err := json.Unmarshal([]byte(data), &registerRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := validator.New().Struct(registerRequest); err != nil {
		ctx.Error(err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file == nil {
			ctx.Error(apperror.FileNotAttachedError())
			return
		}

		ctx.Error(err)
		return
	}

	if err = h.partnerUsecase.AddOnePartner(ctx.Request.Context(), registerRequest, file, *fileHeader); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *PartnerHandler) GetAllPartners(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	partners, err := h.partnerUsecase.GetAllPartners(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, partners)
}

func (h *PartnerHandler) UpdatePartner(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var updateAccountRequest dto.UpdateAccountRequest

	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.EmptyDataError())
		return
	}

	if err := json.Unmarshal([]byte(data), &updateAccountRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := validator.New().Struct(updateAccountRequest); err != nil {
		ctx.Error(err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		if file != nil {
			ctx.Error(err)
			return
		}
	}

	paramPharmacyManagerId := ctx.Param(appconstant.PharmacyManagerIdString)
	pharmacyManagerId, err := strconv.Atoi(paramPharmacyManagerId)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err = h.partnerUsecase.UpdateOnePartner(ctx, updateAccountRequest, int64(pharmacyManagerId), file, fileHeader); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *PartnerHandler) DeletePartner(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	paramPharmacyManagerId := ctx.Param(appconstant.PharmacyManagerIdString)
	pharmacyManagerId, _ := strconv.Atoi(paramPharmacyManagerId)

	if err := h.partnerUsecase.DeleteOnePartner(ctx, int64(pharmacyManagerId)); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *PartnerHandler) SendCredentialsEmail(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var sendEmailRequest dto.SendEmailRequest

	if err := ctx.ShouldBindJSON(&sendEmailRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.partnerUsecase.SendCredentials(ctx.Request.Context(), sendEmailRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}
