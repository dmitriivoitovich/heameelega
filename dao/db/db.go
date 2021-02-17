package db

import (
	"fmt"
	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	connAttemptsDelay = time.Second * 5
)

var (
	DB *gorm.DB
)

func InitConn(conf config.DBConf, logger echo.Logger) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.DBName,
	)

	for {
		db, err := connect(dsn)
		if err != nil {
			logger.Error(err)
			time.Sleep(connAttemptsDelay)

			continue
		}

		DB = db

		break
	}
}

func connect(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&Customer{}); err != nil {
		return nil, err
	}

	return conn, nil
}
