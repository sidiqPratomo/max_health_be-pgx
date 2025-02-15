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

type DrugHandler struct {
	drugUsecase usecase.DrugUsecase
}

func NewDrugHandler(drugUsecase usecase.DrugUsecase) DrugHandler {
	return DrugHandler{
		drugUsecase: drugUsecase,
	}
}

func (h *DrugHandler) GetPharmacyDrugByDrugId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	drugIdStr := ctx.Param(appconstant.DrugIdString)
	latitude := ctx.Query(appconstant.Latitude)
	longitude := ctx.Query(appconstant.Longitude)
	page := ctx.Query(appconstant.Page)
	limit := ctx.Query(appconstant.Limit)

	drugId, err := strconv.Atoi(drugIdStr)
	if err != nil {
		ctx.Error(apperror.DrugIdInvalidError())
		return
	}

	pharmacyDrug, err := h.drugUsecase.GetPharmacyDrugByDrugId(ctx.Request.Context(), int64(drugId), latitude, longitude, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, *pharmacyDrug)
}

func (h *DrugHandler) GetAllDrugsForListing(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	query := util.GetProductQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	resQuery, err := util.ValidateGetProductQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	productList, pageInfo, err := h.drugUsecase.GetAllDrugsForListing(ctx, resQuery)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dto.DrugListingResponse{Drugs: productList, PageInfo: *pageInfo}

	util.ResponseOK(ctx, res)
}

func (h *DrugHandler) UpdateOneDrug(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	drugIdStr := ctx.Param(appconstant.DrugIdString)

	drugId, err := strconv.Atoi(drugIdStr)
	if err != nil {
		ctx.Error(apperror.DrugIdInvalidError())
		return
	}

	var drugRequest dto.UpdateDrugRequest

	dataStr := ctx.Request.FormValue("data")

	err = json.Unmarshal([]byte(dataStr), &drugRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(drugRequest)
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

	err = h.drugUsecase.UpdateOneDrug(ctx.Request.Context(), int64(drugId), drugRequest, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DrugHandler) GetAllDrugs(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	query := util.GetDrugsAdminQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetDrugAdminQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := h.drugUsecase.GetAllDrugs(ctx, validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, res)
}

func (h *DrugHandler) GetDrugByDrugId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	drugIdStr := ctx.Param(appconstant.DrugIdString)

	drugId, err := strconv.Atoi(drugIdStr)
	if err != nil {
		ctx.Error(apperror.DrugIdInvalidError())
		return
	}

	drug, err := h.drugUsecase.GetOneDrugByDrugId(ctx.Request.Context(), int64(drugId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, *drug)
}

func (h *DrugHandler) CreateOneDrug(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var createDrugRequest dto.CreateDrugRequest

	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.EmptyDataError())
		return
	}

	if err := json.Unmarshal([]byte(data), &createDrugRequest); err != nil {
		ctx.Error(err)
		return
	}

	if err := validator.New().Struct(createDrugRequest); err != nil {
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

	if err = h.drugUsecase.CreateOneDrug(ctx.Request.Context(), createDrugRequest, file, fileHeader); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}

func (h *DrugHandler) DeleteOneDrug(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	paramDrugId := ctx.Param(appconstant.DrugIdString)
	drugId, _ := strconv.Atoi(paramDrugId)
	if err := h.drugUsecase.DeleteOneDrug(ctx.Request.Context(), int64(drugId)); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DrugHandler) GetDrugsByPharmacyId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	_, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramPharmacyId := ctx.Param(appconstant.PharmacyIdString)

	queryParams := util.QueryParam{
		Search: ctx.DefaultQuery("search",""),
		Page:  ctx.DefaultQuery("page", ""),
		Limit: ctx.DefaultQuery("limit", ""),
	}

	params, err := util.SetDefaultQueryParams(queryParams)
	if err != nil {
		ctx.Error(err)
		return
	}

	drugsResponse, err := h.drugUsecase.GetDrugsByPharmacyId(ctx, paramPharmacyId, params.Limit, params.Page, params.Search)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, drugsResponse)
}

func (h *DrugHandler) UpdateDrugsByPharmacyDrugId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	_, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	paramPharmacyDrugId := ctx.Param(appconstant.PharmacyDrugIdString)
	updateDrugReq := dto.UpdatePharmacyDrugReq{}
	if err := ctx.ShouldBindJSON(&updateDrugReq); err != nil {
		ctx.Error(err)
		return
	}

	paramPharmacyDrugIdInt, err := strconv.Atoi(paramPharmacyDrugId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.drugUsecase.UpdateDrugsByPharmacyDrugId(ctx, int64(paramPharmacyDrugIdInt), updateDrugReq.Stock, updateDrugReq.Price)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DrugHandler) AddDrugsByPharmacyManager(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	_, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	addDrugReq := dto.AddPharmacyDrugReq{}
	if err := ctx.ShouldBindJSON(&addDrugReq); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.drugUsecase.AddDrugsByPharmacyDrugId(ctx, int64(addDrugReq.PharmacyId), int64(addDrugReq.DrugId), addDrugReq.Stock, addDrugReq.Price); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DrugHandler) DeleteDrugsByPharmacyDrugId(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	paramPharmacyDrugId := ctx.Param(appconstant.PharmacyDrugIdString)
	pharmacyDrugId, _ := strconv.Atoi(paramPharmacyDrugId)

	if err := h.drugUsecase.DeleteDrugsByPharmacyDrugId(ctx, int64(pharmacyDrugId)); err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, nil)
}

func (h *DrugHandler) GetPossibleStockMutation(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	paramPharmacyDrugId := ctx.Param(appconstant.PharmacyDrugIdString)
	pharmacyDrugId, _ := strconv.Atoi(paramPharmacyDrugId)

	pharmacyDrugs, err := h.drugUsecase.GetPossibleStockMutation(ctx, int64(pharmacyDrugId))
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, pharmacyDrugs)
}

func (h *DrugHandler) PostStockMutation(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	paramPharmacyDrugId := ctx.Param(appconstant.PharmacyDrugIdString)
	pharmacyDrugId, _ := strconv.Atoi(paramPharmacyDrugId)
	postStockMutationReq := dto.PostStockMutationRequest{}
	if err := ctx.ShouldBindJSON(&postStockMutationReq); err != nil {
		ctx.Error(err)
		return
	}

	postStockMutationReq.RecipientPharmacyDrugId = int64(pharmacyDrugId)
	err := h.drugUsecase.PostStockMutation(ctx, postStockMutationReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseCreated(ctx, nil)
}
