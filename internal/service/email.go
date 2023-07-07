package service

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"mailService/internal/models"
	"mailService/internal/repository"
	"os"
	"path/filepath"
)

type Email interface {
	AddUser(mail string) (models.Email, error)
	CheckUserByKeyword(keyword string) ([]byte, error) // check if user exists
}

type EmailService struct {
	repo repository.Email
}

func NewEmailService(repo repository.Email) *EmailService {
	return &EmailService{
		repo: repo,
	}
}

func (s *EmailService) AddUser(mail string) (models.Email, error) {
	model, err := s.repo.AddUser(mail)
	if err != nil {
		return models.Email{}, err
	}
	pathFile := "./tmp/" + model.Email
	err = os.MkdirAll(pathFile, 0755)
	if err != nil {
		// TODO LOGG
		return models.Email{}, err
	}

	// Создаем файлы внутри папки
	incomingFile, err := os.Create(filepath.Join(pathFile, "Incoming.txt"))
	if err != nil {
		// TODO LOGG
		return models.Email{}, err
	}
	defer incomingFile.Close()

	outgoingFile, err := os.Create(filepath.Join(pathFile, "Outgoing.txt"))
	if err != nil {
		// TODO LOGG
		return models.Email{}, err
	}
	defer outgoingFile.Close()

	return model, nil
}

func (s *EmailService) CheckUserByKeyword(keyword string) ([]byte, error) {
	model, err := s.repo.CheckUserByKeyword(keyword)
	if err != nil {
		// TODO logg
		return nil, err
	}

	// Путь к исходной папке
	sourceDir := "./tmp/" + model.Email // TODO получать путь из json-конфига
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
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
	log.Println("success")

	return buf.Bytes(), nil
}

// copyFile копирует файл src в dest.
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
