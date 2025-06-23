package configs

import (
	"os"
)

const (
	// Router environment paths
	servicePortEnvPath = "PORT"

	// DB environment paths
	dbHostEnvPath     = "POSTGRES_HOST"
	dbUserEnvPath     = "POSTGRES_USER"
	dbPasswordEnvPath = "POSTGRES_PASSWORD"
	dbNameEnvPath     = "POSTGRES_DB"
	dbPortEnvPath     = "POSTGRES_PORT"

	// Redis environment paths
	redisAddress  = "REDIS_ADDRESS"
	redisPassword = "REDIS_PASSWORD"

	// Notification client environment path
	notificationClientBasePath = "NOTIFICATION_CLIENT_BASE_PATH"
)

type Configs struct {
	Router             *Router
	DB                 *DB
	RedisConfig        *RedisConfig
	NotificationClient *NotificationClient
}

type Router struct {
	Port string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type NotificationClient struct {
	ServiceBaseURL string
}

type DB struct {
	Host     string
	User     string
	Password string
	DB       string
	Port     string
}

func GetConfigs() (*Configs, error) {
	return &Configs{
		Router: &Router{
			Port: os.Getenv(servicePortEnvPath),
		},
		DB: &DB{
			Host:     os.Getenv(dbHostEnvPath),
			User:     os.Getenv(dbUserEnvPath),
			Password: os.Getenv(dbPasswordEnvPath),
			DB:       os.Getenv(dbNameEnvPath),
			Port:     os.Getenv(dbPortEnvPath),
		},
		RedisConfig: &RedisConfig{
			Addr:     os.Getenv(redisAddress),
			Password: os.Getenv(redisPassword),
		},
		NotificationClient: &NotificationClient{
			ServiceBaseURL: os.Getenv(notificationClientBasePath),
		},
	}, nil
}
