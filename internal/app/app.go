package app

import (
	"github.com/sirupsen/logrus"
	"mailService/internal/config"
	"mailService/internal/delivery"
	"mailService/internal/repository"
	server2 "mailService/internal/server"
	service2 "mailService/internal/service"
	"mailService/pkg/database/postgres"
)

func Run(configPath string) {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		logrus.Panic(err)
	}

	db, err := postgres.NewClient(cfg.DbDsn)

	repo := repository.NewEmailRepo(db)

	emailService := service2.NewEmailService(repo, cfg.UserEmailFilesPath)

	handler := delivery.NewHandler(emailService)

	server := server2.NewHTTPServer(cfg.ServerAddr, handler.Init())

	err = server.Run()
	if err != nil {
		logrus.Panic(err)
	}
}
