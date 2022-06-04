package dto

import "github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/domain"

type TweetDto struct {
	Id     int    `json:"id,omitempty"`
	UserId int    `json:"userId" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

func ToEntity(dto *TweetDto) *domain.Tweet {
	return &domain.Tweet{
		Id:     dto.Id,
		UserId: dto.UserId,
		Text:   dto.Text,
	}
}

func ToDto(entity *domain.Tweet) *TweetDto {
	return &TweetDto{
		Id:     entity.Id,
		UserId: entity.UserId,
		Text:   entity.Text,
	}
}
