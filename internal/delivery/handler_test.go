package delivery

import (
	"fmt"
	"log"
	"mailService/internal/repository"
	server2 "mailService/internal/server"
	"mailService/internal/service"
	"mailService/pkg/database/postgres"
	"testing"
)

func TestHandler_Init(t *testing.T) {
	db, err := postgres.NewClient("host=localhost port=5432 user=postgres password=password dbname=mails sslmode=disable timezone=UTC connect_timeout=5")
	if err != nil {
		log.Println(err)
		return
	}
	repo := repository.NewEmailRepo(db)
	serv := service.NewService(repo)
	handler := NewHandler(serv)
	server := server2.NewHTTPServer(":8080", handler.Init())
	fmt.Println("server is started")
	err = server.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
