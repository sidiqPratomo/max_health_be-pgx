package handler

import (
	"github.com/sidiqPratomo/max-health-backend/apperror"
	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	c.Error(apperror.NotFoundError())
}
