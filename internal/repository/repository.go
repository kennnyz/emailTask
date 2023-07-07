package repository

import (
	"database/sql"
	"mailService/internal/service"
)

type Repository struct {
	Email service.EmailRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Email: NewEmailRepo(db),
	}
}
