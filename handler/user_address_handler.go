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

type UserAddressHandler struct {
	userAddressUsecase usecase.UserAddressUsecase
}

func NewUserAddressHandler(userAddressUsecase usecase.UserAddressUsecase) UserAddressHandler {
	return UserAddressHandler{
		userAddressUsecase: userAddressUsecase,
	}
}

func (h *UserAddressHandler) UpdateUserAddress(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var updateUserAddressRequest dto.UpdateUserAddressRequest

	err := ctx.ShouldBindJSON(&updateUserAddressRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	addressId := ctx.Param(appconstant.AddressIdString)

	c := context.WithValue(ctx.Request.Context(), appconstant.AddressIdKey, addressId)

	err = h.userAddressUsecase.UpdateUserAddress(c, accountId.(int64), updateUserAddressRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *UserAddressHandler) GetAllUserAddress(ctx *gin.Context) {
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	address, err := h.userAddressUsecase.GetAllUserAddress(ctx.Request.Context(), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, address)
}

func (h *UserAddressHandler) AddUserAddress(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var addUserAddressRequest dto.AddUserAddressRequest

	err := ctx.ShouldBindJSON(&addUserAddressRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	err = h.userAddressUsecase.AddUserAddress(ctx.Request.Context(), accountId.(int64), addUserAddressRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *UserAddressHandler) AddUserAddressAutofill(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var addUserAddressAutofillRequest dto.AddUserAddressAutofillRequest

	if err := ctx.ShouldBindJSON(&addUserAddressAutofillRequest); err != nil {
		ctx.Error(err)
		return
	}

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	if err := h.userAddressUsecase.AddUserAddressAutofill(ctx.Request.Context(), accountId.(int64), addUserAddressAutofillRequest); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *UserAddressHandler) DeleteUserAddress(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramAddressId := ctx.Param(appconstant.AddressIdString)
	addressId, _ := strconv.Atoi(paramAddressId)
	err := h.userAddressUsecase.DeleteUserAddress(ctx.Request.Context(), int64(addressId), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}
