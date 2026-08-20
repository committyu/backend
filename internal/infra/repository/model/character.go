package model

import(
	"time"
)

type Character struct {
	ID     string
	Name   string
	Job    string
	Hp     int
	Atk    int
	Matk   int
	Def    int
	Mdef   int
	Agi    int
	Luk    int
	Xp     int
	UserID string
	CreatedAt time.Time
}