package models

import "time"

type Email struct {
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UniqueCode string    `json:"unique_code"`
}

type Zip struct {
	Body []byte
	Name string
}
