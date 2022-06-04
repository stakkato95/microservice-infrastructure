package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-go-lib/logger"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/middle/dto"
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
	headers := http.Header{}

	for key, val := range ctx.Request.Header {
		if strings.ToLower(key)[0] == 'x' {
			headers.Set(key, val[0])
		}
	}
	logger.Info(fmt.Sprintf("%v", headers))

	tHeaders := telemetryHeaders{}

	if err := ctx.ShouldBindHeader(&tHeaders); err != nil {
		ctx.JSON(http.StatusOK, err)
	}

	req, err := http.NewRequest("GET", "http://backend-1/request", nil)
	if err != nil {
		errorResponse(ctx, err)
		return
	}
	req.Header = headers

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
		X_request_id: tHeaders.X_request_id,
		Nested:       nestedResponseDto,
	})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
