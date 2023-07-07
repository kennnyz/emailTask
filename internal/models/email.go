package models

type Email struct {
	Email      string `json:"email"`
	UniqueCode string `json:"unique_code"`
}
