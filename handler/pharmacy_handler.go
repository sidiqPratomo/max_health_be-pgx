package handler

import (
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PharmacyHandler struct {
	pharmacyUsecase usecase.PharmacyUsecase
}

func NewPharmacyHandler(pharmacyUsecase usecase.PharmacyUsecase) PharmacyHandler {
	return PharmacyHandler{
		pharmacyUsecase: pharmacyUsecase,
	}
}

func (h *PharmacyHandler) GetPharmacyByManagerId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	queryParams := util.QueryParam{
		Search: ctx.DefaultQuery("search", ""),
		Page:   ctx.DefaultQuery("page", ""),
		Limit:  ctx.DefaultQuery("limit", ""),
	}

	params, err := util.SetDefaultQueryParams(queryParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	partners, err := h.pharmacyUsecase.GetAllPharmacyByManagerId(ctx, accountId.(int64), params.Limit, params.Page, params.Search)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, partners)
}

func (h *PharmacyHandler) AdminGetPharmacyByManagerId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	_, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramManagerId := ctx.Param(appconstant.PharmacyManagerIdString)
	paramManagerIdInt, err := strconv.Atoi(paramManagerId)
	if err != nil {
		ctx.Error(err)
		return
	}

	queryParams := util.QueryParam{
		Search: ctx.DefaultQuery("search", ""),
		Page:   ctx.DefaultQuery("page", ""),
		Limit:  ctx.DefaultQuery("limit", ""),
	}

	params, err := util.SetDefaultQueryParams(queryParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	partners, err := h.pharmacyUsecase.AdminGetAllPharmacyByManagerId(ctx, int64(paramManagerIdInt), params.Limit, params.Page, params.Search)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, partners)
}

func (h *PharmacyHandler) CreateOnePharmacy(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var pharmacyRequest dto.PharmacyRequest

	if err := ctx.ShouldBindJSON(&pharmacyRequest); err != nil {
		ctx.Error(err)
		return
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	if err := h.pharmacyUsecase.CreateOnePharmacy(ctx.Request.Context(), accountId.(int64), pharmacyRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *PharmacyHandler) UpdateOnePharmacy(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var updatePharmacyRequest dto.UpdatePharmacyRequest

	if err := ctx.ShouldBindJSON(&updatePharmacyRequest); err != nil {
		ctx.Error(err)
		return
	}

	for _, pharmacyOperational := range updatePharmacyRequest.Operationals {
		if err := validator.New().Struct(pharmacyOperational); err != nil {
			ctx.Error(err)
			return
		}
	}

	for _, pharmacyCourier := range updatePharmacyRequest.Couriers {
		if err := validator.New().Struct(pharmacyCourier); err != nil {
			ctx.Error(err)
			return
		}
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	pharmacyIdStr := ctx.Param(appconstant.PharmacyIdString)
	pharmacyId, err := strconv.Atoi(pharmacyIdStr)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatePharmacyRequest.Id = int64(pharmacyId)

	if err = h.pharmacyUsecase.UpdateOnePharmacy(ctx, accountId.(int64), updatePharmacyRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *PharmacyHandler) DeleteOnePharmacy(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	pharmacyIdStr := ctx.Param(appconstant.PharmacyIdString)
	pharmacyId, err := strconv.Atoi(pharmacyIdStr)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := h.pharmacyUsecase.DeleteOnePharmacyById(ctx, accountId.(int64), int64(pharmacyId)); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}
