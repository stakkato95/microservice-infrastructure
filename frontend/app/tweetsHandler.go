package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/dto"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/service"
)

type getTweetsUriParams struct {
	UserId int `uri:"userId" binding:"required,min=1"`
}

type TweetsHandler struct {
	service service.TweetsService
}

func (h *TweetsHandler) addTweet(ctx *gin.Context) {
	var tweetDto dto.TweetDto
	if err := ctx.ShouldBindJSON(&tweetDto); err != nil {
		errorResponse(ctx, err)
		return
	}

	createdTweet := h.service.AddTweet(tweetDto)
	ctx.JSON(http.StatusOK, dto.ResponseDto{Data: *createdTweet})
}

func (h *TweetsHandler) getTweets(ctx *gin.Context) {
	var uriParams getTweetsUriParams
	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		errorResponse(ctx, err)
		return
	}

	tweets := h.service.GetAllTweets(uriParams.UserId)
	ctx.JSON(http.StatusOK, dto.ResponseDto{Data: tweets})
}

func errorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseDto{Error: err.Error()})
}
