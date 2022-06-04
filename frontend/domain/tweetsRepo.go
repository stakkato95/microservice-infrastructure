package domain

import (
	"fmt"

	"github.com/stakkato95/service-engineering-go-lib/logger"
	"github.com/stakkato95/service-engineering-microservice-infrastructure/frontend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TweetsRepo interface {
	AddTweet(Tweet) *Tweet
	GetAllTweets(int) []Tweet
}

type postgresTweetsRepo struct {
	db *gorm.DB
}

func NewTweetsRepo() TweetsRepo {
	db, err := gorm.Open(postgres.Open(config.AppConfig.DbSource), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database: " + err.Error())
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Tweet{}); err != nil {
		logger.Fatal("failed to migrate database: " + err.Error())
	}

	return &postgresTweetsRepo{db}
}

func (r *postgresTweetsRepo) AddTweet(tweet Tweet) *Tweet {
	r.db.Create(&tweet)
	logger.Info(fmt.Sprintf("added tweet with id: %d", tweet.Id))
	return &tweet
}

func (r *postgresTweetsRepo) GetAllTweets(userId int) []Tweet {
	tweets := []Tweet{}
	r.db.Where("user_id = ?", userId).Find(&tweets)
	return tweets
}
