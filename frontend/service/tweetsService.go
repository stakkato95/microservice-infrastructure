package service

import (
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/domain"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/dto"
)

type TweetsService interface {
	AddTweet(dto.TweetDto) *dto.TweetDto
	GetAllTweets(int) []dto.TweetDto
}

type defaultTweetsService struct {
	repo domain.TweetsRepo
}

func NewTweetsService(repo domain.TweetsRepo) TweetsService {
	return &defaultTweetsService{repo}
}

func (s *defaultTweetsService) AddTweet(tweetDto dto.TweetDto) *dto.TweetDto {
	entity := dto.ToEntity(&tweetDto)
	createdTweet := s.repo.AddTweet(*entity)
	return dto.ToDto(createdTweet)
}

func (s *defaultTweetsService) GetAllTweets(userId int) []dto.TweetDto {
	tweets := s.repo.GetAllTweets(userId)
	tweetsDto := make([]dto.TweetDto, len(tweets))

	for i, tweet := range tweets {
		tweetsDto[i] = *dto.ToDto(&tweet)
	}

	return tweetsDto
}
