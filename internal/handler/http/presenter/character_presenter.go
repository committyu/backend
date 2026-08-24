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