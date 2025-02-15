package util

import (
	"net/http"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/sidiqPratomo/max-health-backend/dto"
	"github.com/gin-gonic/gin"
)

func ResponseOK(ctx *gin.Context, res any) {
	ctx.JSON(http.StatusOK, dto.Response{Message: appconstant.MsgOK, Data: res})
}

func ResponseCreated(ctx *gin.Context, res any) {
	ctx.JSON(http.StatusCreated, dto.Response{Message: appconstant.MsgCreated, Data: res})
}
