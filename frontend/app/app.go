package app

import (
	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/config"
)

func Start() {
	h := Handler{}

	router := gin.Default()
	router.GET("/request", h.getRequest)
	router.Run(config.AppConfig.ServerPort)
}
