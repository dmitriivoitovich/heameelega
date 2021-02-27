package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	connAttemptsDelay = time.Second * 5
)

var DB *gorm.DB

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
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)

	gormConfig := &gorm.Config{Logger: newLogger}

	// connect to the database
	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// run migrations
	if err := migrate(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	m.InitSchema(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&User{}, &Customer{}); err != nil {
			return err
		}

		if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
			return err
		}

		createSearchIndexQuery := "CREATE INDEX idx_gin_customers_first_name_last_name_lower " +
			"ON customers USING gin (lower(first_name) gin_trgm_ops, lower(last_name) gin_trgm_ops)"

		return tx.Exec(createSearchIndexQuery).Error
	})

	return m.Migrate()
}
