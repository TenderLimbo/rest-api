package main

import (
	"context"
	restapi "github.com/TenderLimbo/rest-api"
	handler "github.com/TenderLimbo/rest-api/pkg/handler"
	"github.com/TenderLimbo/rest-api/pkg/repository"
	service "github.com/TenderLimbo/rest-api/pkg/service"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	var err error
	if err = InitConfig(); err != nil {
		log.Fatalf("failed to init config : %s", err.Error())
	}
	if err = godotenv.Load(); err != nil {
		log.Fatalf("failed to init .env : %s", err.Error())
	}

	db, err := repository.NewPostgresDB(viper.GetStringMapString("db"), os.Getenv("db_password"))
	if err != nil {
		log.Fatalf("failed to connect database : %s", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(restapi.Server)
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Println("listen: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
