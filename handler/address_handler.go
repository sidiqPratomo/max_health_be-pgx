package handler

import (
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	AddressUsecase usecase.AddressUsecase
}

func NewAddressHandler(AddressUsecase usecase.AddressUsecase) AddressHandler {
	return AddressHandler{
		AddressUsecase: AddressUsecase,
	}
}

func (h *AddressHandler) GetAllProvinces(ctx *gin.Context) {
	provinces, err := h.AddressUsecase.GetAllProvinces(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, provinces)
}

func (h *AddressHandler) GetAllCitiesByProvinceCode(ctx *gin.Context) {
	var addressQuery dto.AddressQuery

	if err := ctx.ShouldBindQuery(&addressQuery); err != nil {
		ctx.Error(err)
		return
	}
	cities, err := h.AddressUsecase.GetAllCitiesByProvinceCode(ctx.Request.Context(), addressQuery.ProvinceCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, cities)
}

func (h *AddressHandler) GetAllDistrictsByCityCode(ctx *gin.Context) {
	var addressQuery dto.AddressQuery

	if err := ctx.ShouldBindQuery(&addressQuery); err != nil {
		ctx.Error(err)
		return
	}
	districts, err := h.AddressUsecase.GetAllDistrictByCityCode(ctx.Request.Context(), addressQuery.CityCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, districts)
}

func (h *AddressHandler) GetAllSubdistrictsByDistrictCode(ctx *gin.Context) {
	var addressQuery dto.AddressQuery

	if err := ctx.ShouldBindQuery(&addressQuery); err != nil {
		ctx.Error(err)
		return
	}
	subdistricts, err := h.AddressUsecase.GetAllSubdistrictByDistrictCode(ctx.Request.Context(), addressQuery.DistrictCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, subdistricts)
}
