package handler

import (
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderPharmacyHandler struct {
	orderPharmacyUsecase usecase.OrderPharmacyUsecase
}

func NewOrderPharmacyHandler(orderPharmacyUsecase usecase.OrderPharmacyUsecase) OrderPharmacyHandler {
	return OrderPharmacyHandler{
		orderPharmacyUsecase: orderPharmacyUsecase,
	}
}

func (h *OrderPharmacyHandler) GetOrderPharmacyById(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	orderPharmacyIdStr := ctx.Param(appconstant.OrderPharmacyIdString)

	orderPharmacyId, err := strconv.Atoi(orderPharmacyIdStr)
	if err != nil {
		ctx.Error(apperror.InvalidOrderError())
		return
	}

	orderPharmacyResponse, err := h.orderPharmacyUsecase.GetOneOrderPharmacyById(ctx.Request.Context(), int64(orderPharmacyId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, orderPharmacyResponse)
}

func (h *OrderPharmacyHandler) GetAllOrderPharmacies(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	query := util.GetOrderQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetOrderQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	orderPharmaciesResponse, err := h.orderPharmacyUsecase.GetAllOrderPharmacies(ctx.Request.Context(), validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, orderPharmaciesResponse)
}

func (h *OrderPharmacyHandler) GetAllUserOrderPharmacies(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	query := util.GetOrderQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetOrderQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}
	orderPharmaciesResponse, err := h.orderPharmacyUsecase.GetAllUserOrderPharmacies(ctx.Request.Context(), accountId.(int64), validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, orderPharmaciesResponse)
}

func (h *OrderPharmacyHandler) GetAllPartnerOrderPharmacies(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	query := util.GetOrderQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetOrderQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}
	orderPharmaciesResponse, err := h.orderPharmacyUsecase.GetAllPartnerOrderPharmacies(ctx.Request.Context(), accountId.(int64), validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, orderPharmaciesResponse)
}

func (h *OrderPharmacyHandler) GetAllPartnerOrderPharmaciesSummary(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	orderPharmaciesResponse, err := h.orderPharmacyUsecase.GetAllPartnerOrderPharmaciesSummary(ctx.Request.Context(), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, orderPharmaciesResponse)
}

func (h *OrderPharmacyHandler) UpdateStatusToSent(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramOrderPharmacyId := ctx.Param(appconstant.OrderPharmacyIdString)
	orderPharmacyId, err := strconv.Atoi(paramOrderPharmacyId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.orderPharmacyUsecase.UpdateStatusToSent(ctx, accountId.(int64), int64(orderPharmacyId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *OrderPharmacyHandler) UpdateStatusToConfirmed(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramOrderPharmacyId := ctx.Param(appconstant.OrderPharmacyIdString)
	orderPharmacyId, err := strconv.Atoi(paramOrderPharmacyId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.orderPharmacyUsecase.UpdateStatusToConfirmed(ctx, accountId.(int64), int64(orderPharmacyId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *OrderPharmacyHandler) UpdateStatusToCancelled(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramOrderPharmacyId := ctx.Param(appconstant.OrderPharmacyIdString)
	orderPharmacyId, err := strconv.Atoi(paramOrderPharmacyId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.orderPharmacyUsecase.UpdateStatusToCancelled(ctx, accountId.(int64), int64(orderPharmacyId))
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, nil)
}
