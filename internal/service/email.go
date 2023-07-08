package service

import (
	"archive/zip"
	"bytes"
	"io"
	"mailService/internal/models"
	"os"
	"path/filepath"
	"time"
)

type EmailRepository interface {
	AddUser(mail string) (models.Email, error)
	CheckUserByKeyword(keyword string) (models.Email, error) // check if user exists
	DeleteUsersByTTL(ttl int) error
	DeleteUserByKeyword(keyword string) error
}

type EmailService struct {
	repo         EmailRepository
	usersTmpPath string
	userMailTTl  int
}

func NewEmailService(repo EmailRepository, usersTmpPath string, ttl int) *EmailService {
	return &EmailService{
		repo:         repo,
		usersTmpPath: usersTmpPath,
		userMailTTl:  ttl,
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

func (s *EmailService) GetUserMailZip(keyword string) (models.Zip, error) {
	model, err := s.repo.CheckUserByKeyword(keyword)
	if err != nil {
		return models.Zip{}, err
	}

	expirationTime := model.CreatedAt.Add(time.Duration(s.userMailTTl) * time.Minute)
	if time.Now().After(expirationTime) {
		err := s.deleteUserByKeyword(keyword)
		if err != nil {
			return models.Zip{}, err
		}
		return models.Zip{}, models.NotFoundUserErr
	}
	return s.getUserZip(model)
}

func (s *EmailService) DeleteUsersByTTL() error {
	return s.repo.DeleteUsersByTTL(s.userMailTTl)
}

func (s *EmailService) deleteUserByKeyword(keyword string) error {
	return s.repo.DeleteUserByKeyword(keyword)
}

func (s *EmailService) getUserZip(model models.Email) (models.Zip, error) {
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
		return models.Zip{}, err
	}

	err = zipWriter.Close()
	if err != nil {
		return models.Zip{}, err
	}

	return models.Zip{
		Body: buf.Bytes(),
		Name: model.Email,
	}, nil
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

	pdfFile, err := os.Open("./attachments/readme.pdf")
	if err != nil {
		return err
	}

	defer pdfFile.Close()
	pdfDestination := filepath.Join(pathFile, "readme.pdf")
	pdfDestinationFile, err := os.Create(pdfDestination)
	if err != nil {
		return err
	}
	defer pdfDestinationFile.Close()

	_, err = io.Copy(pdfDestinationFile, pdfFile)
	if err != nil {
		return err
	}

	return nil
}
