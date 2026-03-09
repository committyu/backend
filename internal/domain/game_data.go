package domain

import (
	"time"
)

type GameData struct {
	userID              UserID  
	mainCharacterID     string    
	playTime            int       
	stage               int      
	lastCommitCheckedAt time.Time 
	updatedAt           time.Time
}

func NewGameData(userID UserID) *GameData {
	return &GameData{
		userID:              userID,
		mainCharacterID:     "", 
		playTime:            0,  
		stage:               1,  
		lastCommitCheckedAt: time.Now(),
		updatedAt:           time.Now(),
	}
}

//ゲッター

func (g *GameData) UserID() UserID {
	return g.userID
}

func (g *GameData) MainCharacterID() string {
	return g.mainCharacterID
}

func (g *GameData) PlayTime() int {
	return g.playTime
}

func (g *GameData) Stage() int {
	return g.stage
}

func (g *GameData) LastCommitCheckedAt() time.Time {
	return g.lastCommitCheckedAt
}


func (g *GameData) AddPlayTime(seconds int) {
	if seconds > 0 {
		g.playTime += seconds
		g.updatedAt = time.Now()
	}
}

func (g *GameData) SetLastCommitCheckedAt(t time.Time) {
	g.lastCommitCheckedAt = t
	g.updatedAt = time.Now()
}

func ReconstructGameData(
	userID UserID,
	mainCharacterID string,
	playTime int,
	stage int,
	lastCommitCheckedAt time.Time,
	updatedAt time.Time,
) *GameData {
	return &GameData{
		userID:              userID,
		mainCharacterID:     mainCharacterID,
		playTime:            playTime,
		stage:               stage,
		lastCommitCheckedAt: lastCommitCheckedAt,
		updatedAt:           updatedAt,
	}
}
