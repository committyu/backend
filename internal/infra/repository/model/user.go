package model

import "time"

type User struct {
	ID                string `gorm:"primaryKey"`
	GithubName        string `gorm:"column:github_name"`
	Email             string
	AvatarURL         string
	GithubID          int64  `gorm:"uniqueIndex"`
	GithubAccessToken string `gorm:"type:text"`
	CreatedAt         time.Time
}
