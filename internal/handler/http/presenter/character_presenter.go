package presenter

import "time"

type CreateCharacterRequest struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type CreateCharacterResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Job       string    `json:"job"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type StatusEditCharacterReq struct {
	ID   string `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Xp   int    `json:"xp" validate:"required" minimum:"1" example:"5"`
	Hp   *int   `json:"hp,omitempty" minimum:"1" example:"1"`
	Atk  *int   `json:"atk,omitempty" minimum:"1" example:"2"`
	Matk *int   `json:"matk,omitempty" minimum:"1" example:"1"`
	Def  *int   `json:"def,omitempty" minimum:"1" example:"3"`
	Mdef *int   `json:"mdef,omitempty" minimum:"1" example:"1"`
	Agi  *int   `json:"agi,omitempty" minimum:"1" example:"1"`
	Luk  *int   `json:"luk,omitempty" minimum:"1" example:"1"`
}

type StatusEditCharacterRes struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Job       string    `json:"job"`
	Hp        int       `json:"hp"`
	Atk       int       `json:"atk"`
	Matk      int       `json:"matk"`
	Def       int       `json:"def"`
	Mdef      int       `json:"mdef"`
	Agi       int       `json:"agi"`
	Luk       int       `json:"luk"`
	Xp        int       `json:"xp"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type JobChangeReq struct {
	ID  string `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Job string `json:"job" validate:"required" example:"Hero"`
}
