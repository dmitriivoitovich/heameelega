package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;not null"`
	Email       string     `gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	Password    string     `gorm:"column:password;type:char(60);not null"`
	AccessToken *string    `gorm:"column:access_token;type:uuid;default:null"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp without time zone;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp without time zone;not null"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp without time zone;default:null;index"`
}
