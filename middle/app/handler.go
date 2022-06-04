package app

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/gin-propagate-xheaders/middleware"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/middle/dto"
)

type Handler struct {
}

func (h *Handler) getRequest(ctx *gin.Context) {
	req, err := http.NewRequest("GET", "http://backend-1/request", nil)
	if err != nil {
		errorResponse(ctx, err)
		return
	}
	xheaders := middleware.GetXHeaders(ctx)
	req.Header = xheaders

	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	var nestedResponseDto dto.ServiceResponseDto
	if err := json.NewDecoder(res.Body).Decode(&nestedResponseDto); err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ServiceResponseDto{
		Service:      "middle",
		X_request_id: xheaders["X-Request-Id"][0],
		Nested:       nestedResponseDto,
	})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
