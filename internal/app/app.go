package app

import (
	"github.com/sirupsen/logrus"
	"mailService/internal/config"
	"mailService/internal/delivery"
	"mailService/internal/repository"
	server2 "mailService/internal/server"
	service2 "mailService/internal/service"
	"mailService/pkg/database/postgres"
	"time"
)

func Run(configPath string) {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		logrus.Panic(err)
	}

	db, err := postgres.NewClient(cfg.DbDsn)

	repo := repository.NewEmailRepo(db)

	// delete users by ttl
	emailService := service2.NewEmailService(repo, cfg.UserEmailFilesPath, cfg.TimeToLiveLink)
	ticker := time.NewTimer(time.Minute * time.Duration(cfg.TimeToLiveLink))
	go func() {
		for {
			select {
			case <-ticker.C:
				err := emailService.DeleteUsersByTTL()
				if err != nil {
					logrus.Error(err)
				}
				ticker.Reset(time.Minute * time.Duration(cfg.TimeToLiveLink))
			}
		}
	}()

	handler := delivery.NewHandler(emailService)

	server := server2.NewHTTPServer(cfg.ServerAddr, handler.Init())

	err = server.Run()
	if err != nil {
		logrus.Panic(err)
	}
}
