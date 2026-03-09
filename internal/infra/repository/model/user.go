package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey"`
	Name      string
	Email     string
	AvatarURL string
	GithubID  int64     `gorm:"uniqueIndex"`
	CreatedAt time.Time
}