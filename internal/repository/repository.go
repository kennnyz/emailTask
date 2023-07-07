package repository

import "database/sql"

type Repository struct {
	Email Email
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Email: NewEmailRepo(db),
	}
}
