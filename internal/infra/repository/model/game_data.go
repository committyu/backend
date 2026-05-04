package model

import "time"

type GameData struct {
	UserID              string `gorm:"primaryKey"`
	MainCharacterID     string
	PlayTime            int
	Stage               int
	GithubTotalCommits  int
	LastCommitCheckedAt time.Time
	UpdatedAt           time.Time
}