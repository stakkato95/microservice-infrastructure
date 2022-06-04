package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-go-lib/logger"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/middle/config"
)

const XHeadersKey = "x-headers"

func XHeadersPropagation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headers := http.Header{}

		for key, val := range ctx.Request.Header {
			if strings.ToLower(key)[0] == 'x' {
				headers.Set(key, val[0])
			}
		}

		ctx.Set(XHeadersKey, headers)
		ctx.Next()
	}
}

func GetXHeaders(ctx *gin.Context) http.Header {
	headersRaw, ok := ctx.Get(XHeadersKey)
	if !ok {
		logger.Fatal("XHeadersPropagation is not used")
	}
	headers, ok := headersRaw.(http.Header)
	if !ok {
		logger.Fatal("XHeadersPropagation cast error")
	}
	return headers
}

func Start() {
	h := Handler{}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(XHeadersPropagation())

	router.GET("/request", h.getRequest)
	router.Run(config.AppConfig.ServerPort)
}
