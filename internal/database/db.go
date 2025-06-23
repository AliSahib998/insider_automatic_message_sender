package database

import (
	"insider_task/internal/configs"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectDB(config *configs.DB) (*gorm.DB, error) {
	dsn := "host=" + config.Host +
		" user=" + config.User +
		" password=" + config.Password +
		" dbname=" + config.DB +
		" port=" + config.Port +
		" sslmode=disable" +
		" TimeZone=UTC"

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
