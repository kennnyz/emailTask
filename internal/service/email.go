package service

import (
	"archive/zip"
	"bytes"
	"io"
	"mailService/internal/models"
	"os"
	"path/filepath"
)

type EmailRepository interface {
	AddUser(mail string) (models.Email, error)
	CheckUserByKeyword(keyword string) (models.Email, error) // check if user exists
}

type EmailService struct {
	repo         EmailRepository
	usersTmpPath string
}

func NewEmailService(repo EmailRepository, usersTmpPath string) *EmailService {
	return &EmailService{
		repo:         repo,
		usersTmpPath: usersTmpPath,
	}
}

func (s *EmailService) AddUser(mail string) (models.Email, error) {
	model, err := s.repo.AddUser(mail)
	if err != nil {
		return models.Email{}, err
	}

	err = s.addUserMailFolder(model)
	if err != nil {
		return models.Email{}, err
	}

	return model, nil
}

func (s *EmailService) CheckUserByKeyword(keyword string) ([]byte, error) {
	model, err := s.repo.CheckUserByKeyword(keyword)
	if err != nil {
		return nil, err
	}

	return s.getUserZip(model)
}

func (s *EmailService) getUserZip(model models.Email) ([]byte, error) {
	sourceDir := s.usersTmpPath + model.Email
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Получаем относительный путь файла
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Игнорируем директории
		if info.IsDir() {
			return nil
		}

		// Создаем новый zip-файл в архиве
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Открываем исходный файл
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		// Копируем содержимое файла в zip-файл
		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *EmailService) addUserMailFolder(model models.Email) error {
	pathFile := s.usersTmpPath + model.Email
	err := os.MkdirAll(pathFile, 0755)
	if err != nil {
		return err
	}

	incomingFile, err := os.Create(filepath.Join(pathFile, "Incoming.txt"))
	if err != nil {
		return err
	}
	defer incomingFile.Close()

	outgoingFile, err := os.Create(filepath.Join(pathFile, "Outgoing.txt"))
	if err != nil {
		return err
	}
	defer outgoingFile.Close()

	return nil
}
