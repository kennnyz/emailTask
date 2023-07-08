package repository

import (
	"database/sql"
	"fmt"
	"mailService/internal/models"
)

type EmailRepo struct {
	db *sql.DB
}

func NewEmailRepo(db *sql.DB) *EmailRepo {
	return &EmailRepo{db: db}
}

func (e *EmailRepo) AddUser(mail string) (models.Email, error) {
	var email models.Email

	checkQuery := "SELECT email, unique_code, created_at FROM users WHERE email = $1"
	err := e.db.QueryRow(checkQuery, mail).Scan(&email.Email, &email.UniqueCode, &email.CreatedAt)
	if err == nil {
		return models.Email{}, models.UserAlreadyExistErr
	} else if err != sql.ErrNoRows {
		return models.Email{}, err
	}

	insertQuery := "INSERT INTO users (email) VALUES ($1) RETURNING email, unique_code, created_at"
	err = e.db.QueryRow(insertQuery, mail).Scan(&email.Email, &email.UniqueCode, &email.CreatedAt)
	if err != nil {
		return models.Email{}, err
	}

	return email, nil
}

func (e *EmailRepo) CheckUserByKeyword(keyword string) (models.Email, error) {
	var model models.Email
	checkQuery := "SELECT email, unique_code, created_at FROM users WHERE unique_code = $1"
	err := e.db.QueryRow(checkQuery, keyword).Scan(&model.Email, &model.UniqueCode, &model.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Email{}, models.NotFoundUserErr
	} else if err != nil {
		return models.Email{}, err
	}
	return model, nil
}

func (e *EmailRepo) DeleteUsersByTTL(ttl int) error {
	deleteQuery := fmt.Sprintf("DELETE FROM users WHERE created_at < NOW() - INTERVAL '%d minute'", ttl)
	_, err := e.db.Exec(deleteQuery)
	if err != nil {
		return err
	}

	return nil
}

func (e *EmailRepo) DeleteUserByKeyword(keyword string) error {
	deleteQuery := "DELETE FROM users WHERE unique_code = $1"
	_, err := e.db.Exec(deleteQuery, keyword)
	if err != nil {
		return err
	}

	return nil
}
