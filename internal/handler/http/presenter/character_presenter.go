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

type EditCharacterReq struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
	Hp   *int    `json:"hp"`
	Atk  *int    `json:"atk"`
	Matk *int    `json:"matk"`
	Def  *int    `json:"def"`
	Mdef *int    `json:"mdef"`
	Agi  *int    `json:"agi"`
	Luk  *int    `json:"luk"`
	Xp   *int    `json:"xp"`
}

type EditCharacterRes struct {
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
