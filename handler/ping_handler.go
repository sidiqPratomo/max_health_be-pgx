package handler

import (
	"github.com/sidiqPratomo/max-health-backend/util"
	"github.com/gin-gonic/gin"
)

type PingHandlerOpts struct {
}

type PingHandler struct {
}

func NewPingHandler(opts PingHandlerOpts) *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(ctx *gin.Context) {
	util.ResponseOK(ctx, "pong")
}
