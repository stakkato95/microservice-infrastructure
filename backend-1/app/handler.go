package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/backend-1/dto"
)

type Handler struct {
}

type telemetryHeaders struct {
	X_request_id      string `header:"x-request-id"`
	X_b3_traceid      string `header:"x-b3-traceid"`
	X_b3_spanid       string `header:"x-b3-spanid"`
	X_b3_parentspanid string `header:"x-b3-parentspanid"`
	X_b3_sampled      string `header:"x-b3-sampled"`
	X_b3_flags        string `header:"x-b3-flags"`
	B3                string `header:"b3"`
}

func (h *telemetryHeaders) ToHeaders() http.Header {
	return http.Header{
		"x-request-id":      {h.X_request_id},
		"x-b3-traceid":      {h.X_b3_traceid},
		"x-b3-spanid":       {h.X_b3_spanid},
		"x-b3-parentspanid": {h.X_b3_parentspanid},
		"x-b3-sampled":      {h.X_b3_sampled},
		"x-b3-flags":        {h.X_b3_flags},
		"b3":                {h.B3},
	}
}

func (h *Handler) getRequest(ctx *gin.Context) {
	tHeaders := telemetryHeaders{}

	if err := ctx.ShouldBindHeader(&tHeaders); err != nil {
		ctx.JSON(http.StatusOK, err)
	}

	ctx.JSON(http.StatusOK, dto.ServiceResponseDto{
		Service:      "backend-1",
		X_request_id: tHeaders.X_request_id,
	})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
