package models

import "gorm.io/gorm"

type Mweet struct {
	Image   string `json:"img"`
	Content string `json:"content"`
	UserId  string `json:"userId"`
	gorm.Model
}

type User struct {
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	OAuth_Type string `json:"oAuth_Type"`
	gorm.Model
}

func (Mweet) TableName() string {
	return "mweets"
}
