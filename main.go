package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "insider_task/docs"
	clients "insider_task/internal/client"
	"insider_task/internal/configs"
	"insider_task/internal/database"
	"insider_task/internal/handlers"
	"insider_task/internal/repositories"
	"insider_task/internal/router"
	"insider_task/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Package main ginSwagger
//
// @title Insider Message Service
// @version 1.0
// @description Gin API with Swagger documentation Message service for Insider
// @BasePath /api/v1/messages
func main() {
	// Initialize .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize configs
	appConfigs, err := configs.GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := database.ConnectDB(appConfigs.DB)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		appConfigs.DB.User,
		appConfigs.DB.Password,
		appConfigs.DB.Host,
		appConfigs.DB.Port,
		appConfigs.DB.DB)

	m, err := migrate.New(
		"file://db/migrations",
		connectionString,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	// Connect to redis
	redisClient, err := database.ConnectRedis(context.TODO(), appConfigs.RedisConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	messageRepo := repositories.NewMessagesRepository(db)

	// initialize clients
	appClients := clients.NewClients(appConfigs)

	// Initialize services
	messageService := service.NewMessageService(messageRepo, appClients, redisClient)

	// Initialize handlers
	handler := handlers.NewMessageHandler(messageService)

	appRouter := router.NewRouter(appConfigs.Router, handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	wg, _ := errgroup.WithContext(ctx)

	go func() {
		if err = appRouter.InitAndRun(); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	cancel()

	if err = wg.Wait(); err != nil {
		log.Fatal(err)
	}
}
