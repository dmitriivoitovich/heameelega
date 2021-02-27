package db

import (
	"fmt"
	"time"

	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	connAttemptsDelay = time.Second * 5
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;not null"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp without time zone;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp without time zone;not null"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp without time zone;default:null;index"`
}

var DB *gorm.DB

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	if !m.CreatedAt.IsZero() {
		tx.Statement.SetColumn("CreatedAt", m.CreatedAt)
	}

	if !m.UpdatedAt.IsZero() {
		tx.Statement.SetColumn("UpdatedAt", m.UpdatedAt)
	}

	if m.DeletedAt != nil && !m.DeletedAt.IsZero() {
		tx.Statement.SetColumn("DeletedAt", m.DeletedAt)
	}

	return nil
}

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

	if err := conn.AutoMigrate(&User{}, &Customer{}); err != nil {
		return nil, err
	}

	return conn, nil
}
