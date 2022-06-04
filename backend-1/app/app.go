package app

import (
	"github.com/gin-gonic/gin"
	"github.com/stakkato95/gin-propagate-xheaders/middleware"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/backend-1/config"
)

func Start() {
	h := Handler{}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.XHeadersPropagation())

	router.GET("/request", h.getRequest)
	router.Run(config.AppConfig.ServerPort)
}
