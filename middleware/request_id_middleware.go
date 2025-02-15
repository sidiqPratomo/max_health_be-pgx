package middleware

import (
	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func RequestIdHandlerMiddleware(c *gin.Context) {
	uuid := uuid.NewString()
	c.Set(appconstant.RequestId, uuid)

	c.Next()
}
