package handler

import (
	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/sidiqPratomo/max-health-backend/usecase"
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	StockUsecase usecase.StockUsecase
}

func NewStockHandler(stockUsecase usecase.StockUsecase) StockHandler {
	return StockHandler{
		StockUsecase: stockUsecase,
	}
}

func (h *StockHandler) GetAllStockChanges(ctx *gin.Context) {
	accountId, exists := ctx.Get(appconstant.AccountId)
	if !exists {
		ctx.Error(apperror.UnauthorizedError())
		return
	}

	var stockChangeQuery dto.StockChangeQuery

	if err := ctx.ShouldBindQuery(&stockChangeQuery); err != nil {
		ctx.Error(err)
		return
	}

	stockChanges, err := h.StockUsecase.GetAllStockChanges(ctx.Request.Context(), accountId.(int64), stockChangeQuery.PharmacyId)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseOK(ctx, stockChanges)
}
