package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/dto"
)

type Handler struct {
}

func (h *Handler) getRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.ResponseDto{Data: "data"})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
