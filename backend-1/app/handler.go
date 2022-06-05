package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/gin-propagate-xheaders/middleware"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/backend-1/dto"
)

type Handler struct {
}

func (h *Handler) getRequest(ctx *gin.Context) {
	xheaders := middleware.GetXHeaders(ctx)
	ctx.JSON(http.StatusOK, dto.ServiceResponseDto{
		Service:       "backend-1",
		X_request_id:  xheaders["X-Request-Id"][0],
		X_api_user_id: xheaders["X-Api-User-Id"][0],
	})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
