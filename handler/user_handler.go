package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	profileResponse, err := h.userUsecase.GetProfile(ctx.Request.Context(), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, profileResponse)
}

func (h *UserHandler) UpdateData(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	var updateUserDataRequest dto.UpdateUserDataRequest
	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("data is empty"), "data is empty"))
		return
	}
	err := json.Unmarshal([]byte(data), &updateUserDataRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(updateUserDataRequest)
	if err != nil {
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

	accountData := dto.UpdateUserDataRequestToDetailedUser(updateUserDataRequest)
	id, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	stringId := id.(int64)
	accountData.Id = stringId

	err = h.userUsecase.UpdateData(ctx.Request.Context(), accountData, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)

}
