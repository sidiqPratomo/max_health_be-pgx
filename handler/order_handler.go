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

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) OrderHandler {
	return OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) CheckoutOrder(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var orderCheckoutRequest dto.OrderCheckoutRequest
	if err := ctx.ShouldBindJSON(&orderCheckoutRequest); err != nil {
		ctx.Error(err)
		return
	}

	for _, pharmacy := range orderCheckoutRequest.Pharmacies {
		err := validator.New().Struct(pharmacy)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	orderCheckoutRequest.AccountId = accountId.(int64)

	orderId, err := h.orderUsecase.CheckoutOrder(ctx.Request.Context(), orderCheckoutRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dto.OrderCheckoutResponse{OrderId: *orderId}
	util.ResponseCreated(ctx, res)
}

func (h *OrderHandler) ConfirmPayment(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	orderId, err := strconv.Atoi(ctx.Param(appconstant.OrderIdString))
	if err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if orderId < 1 {
		ctx.Error(apperror.OrderNotFoundError())
		return
	}

	var req dto.OrderChangeStatusRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	if req.StatusId != 1 && req.StatusId != 3 {
		ctx.Error(apperror.InvalidOrderStatusError())
		return
	}

	err = h.orderUsecase.ConfirmPayment(ctx, int64(orderId), req.StatusId)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *OrderHandler) UploadPaymentProofOrder(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
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

	paramOrderId := ctx.Param(appconstant.OrderIdString)
	orderId, err := strconv.Atoi(paramOrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err = h.orderUsecase.UploadPaymentProofOrder(ctx.Request.Context(), accountId.(int64), int64(orderId), file, *fileHeader); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *OrderHandler) GetAllOrders(ctx *gin.Context) {
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

	ordersResponse, err := h.orderUsecase.GetAllOrders(ctx.Request.Context(), validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, ordersResponse)
}

func (h *OrderHandler) GetAllUserPendingOrders(ctx *gin.Context) {
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

	ordesResponse, err := h.orderUsecase.GetAllUserPendingOrders(ctx.Request.Context(), accountId.(int64), validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, ordesResponse)
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramOrderId := ctx.Param(appconstant.OrderIdString)
	orderId, err := strconv.Atoi(paramOrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err = h.orderUsecase.CancelOrder(ctx.Request.Context(), accountId.(int64), int64(orderId)); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *OrderHandler) GetOrderById(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	orderIdStr := ctx.Param(appconstant.OrderIdString)

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		ctx.Error(apperror.InvalidOrderError())
		return
	}

	orderResponse, err := h.orderUsecase.GetOneOrderById(ctx.Request.Context(), int64(orderId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, orderResponse)
}
	

 
