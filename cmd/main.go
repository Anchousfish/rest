package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anchousfish/rest"
	"github.com/Anchousfish/rest/pkg/handler"
	"github.com/Anchousfish/rest/pkg/repository"
	"github.com/Anchousfish/rest/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("fail while reading config file: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("fail to load env variable :%s", err.Error())
	}
	db, err := repository.NewPostgresDb(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(rest.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while trying to start server %s", err.Error())
		}
	}()
	logrus.Println("rest app started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("rest ap is shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while trying to stop the server : %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured while trying to close db connection : %s", err.Error())
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
