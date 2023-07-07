package service

import "mailService/internal/repository"

type Service struct {
	EmailService Email
}

func NewService(repo repository.Email) *Service {
	return &Service{
		EmailService: NewEmailService(repo),
	}
}
