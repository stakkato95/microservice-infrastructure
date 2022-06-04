package app

import (
	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/config"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/domain"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/service"
)

func Start() {
	repo := domain.NewTweetsRepo()
	service := service.NewTweetsService(repo)

	h := TweetsHandler{service}

	router := gin.Default()
	router.POST("/tweets", h.addTweet)
	router.GET("/tweets/:userId", h.getTweets)
	router.Run(config.AppConfig.ServerPort)
}
