package dblayer

import "github.com/mwitter-backend-gRPC/server/models"

type DBLayer interface {
	GetAllMweeter() ([]models.Mweet, error)
	CreateMweet(*models.Mweet) error
	UpdateMweet(mweetId int, mweet *models.Mweet) error
	DeleteMweet(mweetId int) error
	GetMweeterById(mweetId int) (*models.Mweet, error)
	ExistUser(userId int) bool
}
