package domain

import (
	"time"
)

type GameData struct {
	userID              UserID    // どのユーザーのデータか
	mainCharacterID     string    // 現在選択中のキャラID
	playTime            int       // 秒単位の残りプレイ時間
	stage               int       // 現在の到達ステージ
	lastCommitCheckedAt time.Time // 最後にGitHubを確認した時間
	updatedAt           time.Time
}

// NewGameData は新規ユーザー用の初期データを生成します
func NewGameData(userID UserID) *GameData {
	return &GameData{
		userID:              userID,
		mainCharacterID:     "", // 初期状態は空、またはデフォルト値
		playTime:            0,  // 最初は0。コミットで増える
		stage:               1,  // ステージ1から開始
		lastCommitCheckedAt: time.Now(),
		updatedAt:           time.Now(),
	}
}

// --- ゲッター ---

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

// --- ビジネスロジック（状態変更） ---

// AddPlayTime はコミット数に応じて時間を加算します
func (g *GameData) AddPlayTime(seconds int) {
	if seconds > 0 {
		g.playTime += seconds
		g.updatedAt = time.Now()
	}
}

// SetLastCommitCheckedAt はチェック時刻を更新します
func (g *GameData) SetLastCommitCheckedAt(t time.Time) {
	g.lastCommitCheckedAt = t
	g.updatedAt = time.Now()
}