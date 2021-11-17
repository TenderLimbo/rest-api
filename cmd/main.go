package main

import (
	restapi "github.com/TenderLimbo/rest-api"
	handler "github.com/TenderLimbo/rest-api/pkg/handler"
	"github.com/TenderLimbo/rest-api/pkg/repository"
	service "github.com/TenderLimbo/rest-api/pkg/service"
	"log"
)

func main() {
	var err error
	db, err := repository.NewSqliteDB()
	if err != nil {
		log.Fatalf("failed to connect database : %s", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(restapi.Server)
	if err = srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error while running server : %s", err.Error())
	}
}
