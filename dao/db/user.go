package db

import (
	"time"

	"github.com/google/uuid"
)

const (
	LanguageCodeRussian = "RU"
	LanguageCodeEnglish = "EN"
)

type User struct {
	BaseModel
	Email       string     `gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	Password    string     `gorm:"column:password;type:char(60);not null"`
	AccessToken *uuid.UUID `gorm:"column:access_token;type:uuid;default:null"`
	Language    string     `gorm:"column:language;type:char(2);not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp without time zone;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp without time zone;not null"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp without time zone;default:null;index"`
	Customers   []Customer `gorm:"ForeignKey:UserID"`
}
