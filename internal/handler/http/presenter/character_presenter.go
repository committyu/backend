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
	ID   string `json:"id"`
	Xp   int    `json:"xp"`
	Hp   *int   `json:"hp"`
	Atk  *int   `json:"atk"`
	Matk *int   `json:"matk"`
	Def  *int   `json:"def"`
	Mdef *int   `json:"mdef"`
	Agi  *int   `json:"agi"`
	Luk  *int   `json:"luk"`
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
