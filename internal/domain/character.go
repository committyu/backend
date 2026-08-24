package domain

import (
	"time"
)

type Character struct {
	id        CharacterID
	name      string
	job       string
	hp        int
	atk       int
	matk      int
	def       int
	mdef      int
	agi       int
	luk       int
	xp        int
	userID    UserID
	createdAt time.Time
}

func NewCharacter(name string, job string, userID UserID, createdAt time.Time) *Character {
	return &Character{
		id:        NewCharacterID(),
		name:      name,
		job:       job,
		hp:        10,
		atk:       0,
		matk:      0,
		def:       0,
		mdef:      0,
		agi:       0,
		luk:       0,
		xp:        0,
		userID:    userID,
		createdAt: createdAt,
	}
}

// RestoreCharacter rebuilds a Character from persisted values.
func RestoreCharacter(
	id CharacterID,
	name string,
	job string,
	hp, atk, matk, def, mdef, agi, luk, xp int,
	userID UserID,
	createdAt time.Time,
) *Character {
	return &Character{
		id: id, name: name, job: job, hp: hp, atk: atk, matk: matk,
		def: def, mdef: mdef, agi: agi, luk: luk, xp: xp,
		userID: userID, createdAt: createdAt,
	}
}

type CharacterUpdate struct {
	Name *string
	Hp   *int
	Atk  *int
	Matk *int
	Def  *int
	Mdef *int
	Agi  *int
	Luk  *int
	Xp   *int
}

func (c *Character) ID() CharacterID {
	return c.id
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) Job() string {
	return c.job
}

func (c *Character) Hp() int {
	return c.hp
}

func (c *Character) Atk() int {
	return c.atk
}

func (c *Character) Matk() int {
	return c.matk
}

func (c *Character) Def() int {
	return c.def
}

func (c *Character) Mdef() int {
	return c.mdef
}

func (c *Character) Agi() int {
	return c.agi
}

func (c *Character) Luk() int {
	return c.luk
}

func (c *Character) Xp() int {
	return c.xp
}

func (c *Character) UserID() UserID {
	return c.userID
}

func (u *Character) CreatedAt() time.Time {
	return u.createdAt
}
