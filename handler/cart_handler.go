package handler

import (
	"context"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUsecase usecase.CartUsecase
}

func NewCartHandler(cartUsecase usecase.CartUsecase) CartHandler {
	return CartHandler{
		cartUsecase: cartUsecase,
	}
}

func (h *CartHandler) CalculateDeliveryFee(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var deliveryFeeRequest dto.DeliveryFeeRequest
	if err := ctx.ShouldBindJSON(&deliveryFeeRequest); err != nil {
		ctx.Error(err)
		return
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	deliveryFeeRequest.AccountId = accountId.(int64)

	deliveryFees, err := h.cartUsecase.CalculateDeliveryFee(ctx.Request.Context(), deliveryFeeRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, deliveryFees)
}
	
func (h *CartHandler) CreateOneCart(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	value, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	c := context.WithValue(ctx, appconstant.AccountIdKey, value.(int64))

	cartReq := dto.CreateOneCartRequest{}
	if err := ctx.ShouldBindJSON(&cartReq); err != nil {
		ctx.Error(err)
		return
	}

	err := h.cartUsecase.CreateOneCart(c, cartReq.PharmacyDrugId)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseCreated(ctx, nil)
}

func (h *CartHandler) UpdateQtyCart(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	value, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	c := context.WithValue(ctx, appconstant.AccountIdKey, value.(int64))

	cartIdInt, err := strconv.Atoi(ctx.Param(appconstant.CartIdString))
	if err != nil {
		ctx.Error(err)
		return
	}

	if cartIdInt < 1 {
		ctx.Error(err)
		return
	}

	cartReq := dto.UpdateQtyCartRequest{}
	if err := ctx.ShouldBindJSON(&cartReq); err != nil {
		ctx.Error(err)
		return
	}

	err = h.cartUsecase.UpdateOneCart(c, int64(cartIdInt), cartReq.Quantity)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *CartHandler) DeleteOneCart(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	value, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	c := context.WithValue(ctx, appconstant.AccountIdKey, value.(int64))

	cartIdInt, err := strconv.Atoi(ctx.Param(appconstant.CartIdString))
	if err != nil {
		ctx.Error(err)
		return
	}

	if cartIdInt < 1 {
		ctx.Error(apperror.InvalidCartItemError())
		return
	}

	err = h.cartUsecase.DeleteOneCart(c, int64(cartIdInt))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *CartHandler) GetAllCart(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	value, exist := ctx.Get(appconstant.AccountId)
	if !exist {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	c := context.WithValue(ctx, appconstant.AccountIdKey, value.(int64))

	queryParams := util.QueryParam{
		Page:  ctx.DefaultQuery("page", ""),
		Limit: ctx.DefaultQuery("limit", ""),
	}

	params, err := util.SetDefaultQueryParams(queryParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	cart, err := h.cartUsecase.GetAllCartById(c, params.Limit, params.Page)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, cart)
}
