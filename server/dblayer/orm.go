package dblayer

import (
	"fmt"

	"github.com/mwitter-backend-gRPC/server/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBORM struct {
	*gorm.DB
}

func NewORM(dbname string, con gorm.Config) (*DBORM, error) {

	dsn := fmt.Sprintf("root@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=true", dbname)
	dsn = dsn + "&loc=Asia%2FSeoul"
	db, err := gorm.Open(mysql.Open(dsn), &con)

	db.AutoMigrate(&models.Mweet{})

	return &DBORM{
		DB: db,
	}, err
}

func (db *DBORM) GetAllMweeter() (mweets []models.Mweet, err error) {

	result := db.Table("mweets").Select("id", "content", "image", "user_id")

	return mweets, result.Find(&mweets).Error
}

func (db *DBORM) CreateMweet(mweet *models.Mweet) error {

	return db.Create(&mweet).Error
}

func (db *DBORM) UpdateMweet(mweetId int, mweet *models.Mweet) error {
	return db.Table("mweets").Where("id = ?", mweetId).Updates(mweet).Error
}

func (db *DBORM) DeleteMweet(mweetId int) error {

	return db.Where("id = ?", mweetId).Delete(&models.Mweet{}).Error
}

func (db *DBORM) GetMweeterById(mweetId int) (mweet *models.Mweet, err error) {

	result := db.Table("mweets").Where("id = ?", mweetId)
	return mweet, result.Find(&mweet).Error
}

func (db *DBORM) ExistUser(userId int) bool {
	user := &models.User{}
	err := db.Table("users").Where("id = ?", userId).Find(user).Error
	if err != nil {
		return false
	}

	if user.ID == 0 {
		return false
	}

	return true
}
