package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DoctorHandler struct {
	doctorUsecase usecase.DoctorUsecase
}

func NewDoctorHandler(doctorUsecase usecase.DoctorUsecase) DoctorHandler {
	return DoctorHandler{
		doctorUsecase: doctorUsecase,
	}
}

func (h *DoctorHandler) GetAllDoctors(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	queryParams := util.QueryParam{
		Sort:             ctx.DefaultQuery("sort", ""),
		SortBy:           ctx.DefaultQuery("sortBy", ""),
		Page:             ctx.DefaultQuery("page", ""),
		Limit:            ctx.DefaultQuery("limit", ""),
		SpecializationId: ctx.DefaultQuery("specialization", ""),
	}

	params, err := util.SetDefaultQueryParams(queryParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	doctors, err := h.doctorUsecase.GetAllDoctors(ctx.Request.Context(), params.Sort, params.SortBy, params.Limit, params.SpecializationId, params.Page)
	if err != nil {
		ctx.Error(err)
		return
	}
	util.ResponseOK(ctx, doctors)
}

func (h *DoctorHandler) UpdateData(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	var updateDoctorDataRequest dto.UpdateDoctorDataRequest
	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.NewAppError(http.StatusBadRequest, errors.New("data is empty"), "data is empty"))
		return
	}
	err := json.Unmarshal([]byte(data), &updateDoctorDataRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(updateDoctorDataRequest)
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

	accountData := dto.UpdateDoctorDataRequestToDetailedDoctor(updateDoctorDataRequest)
	id, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}
	stringId := id.(int64)
	accountData.Id = stringId

	err = h.doctorUsecase.UpdateData(ctx.Request.Context(), accountData, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DoctorHandler) GetAllDoctorSpecialization(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	specializationList, err := h.doctorUsecase.GetAllDoctorSpecialization(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, specializationList)
}

func (h *DoctorHandler) GetProfile(ctx *gin.Context) {
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	profileResponse, err := h.doctorUsecase.GetProfile(ctx.Request.Context(), accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, profileResponse)
}

func (h *DoctorHandler) GetProfileForPublic(ctx *gin.Context) {
	doctorIdString := ctx.Param(appconstant.DoctorIdString)
	doctorId, _ := strconv.Atoi(doctorIdString)

	profileResponse, err := h.doctorUsecase.GetProfileForPublic(ctx.Request.Context(), int64(doctorId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, profileResponse)
}

func (h *DoctorHandler) UpdateDoctorStatus(ctx *gin.Context) {
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	var request dto.UpdateDoctorStatusRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.doctorUsecase.UpdateDoctorStatus(ctx.Request.Context(), accountId.(int64), request.IsOnline)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DoctorHandler) GetDoctorIsOnline(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	isOnline, err := h.doctorUsecase.GetDoctorIsOnline(ctx, accountId.(int64))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, isOnline)
}
