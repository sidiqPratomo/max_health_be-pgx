package handler

import (
	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(reportUsecase usecase.ReportUsecase) ReportHandler {
	return ReportHandler{
		reportUsecase: reportUsecase,
	}
}

func (h *ReportHandler) GetPharmacyDrugCategoryReport(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	query := util.GetReportQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetReportQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrugCategoryReport, err := h.reportUsecase.GetPharmacyDrugCategoryReport(ctx.Request.Context(), accountId.(int64), *validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, pharmacyDrugCategoryReport)
}

func (h *ReportHandler) GetPharmacyDrugReport(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	query := util.GetReportQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetReportQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrugReport, err := h.reportUsecase.GetPharmacyDrugReport(ctx.Request.Context(), accountId.(int64), *validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, pharmacyDrugReport)
}

func (h *ReportHandler) GetDrugCategoryReport(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	query := util.GetReportQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetReportQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	drugCategoryReport, err := h.reportUsecase.GetDrugCategoryReport(ctx.Request.Context(), *validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, drugCategoryReport)
}

func (h *ReportHandler) GetDrugReport(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	query := util.GetReportQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(apperror.BadRequestError(err))
		return
	}

	if err := validator.New().Struct(query); err != nil {
		ctx.Error(err)
		return
	}

	validatedQuery, err := util.ValidateGetReportQuery(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	drugReport, err := h.reportUsecase.GetDrugReport(ctx.Request.Context(), *validatedQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, drugReport)
}
