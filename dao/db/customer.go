package db

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	GenderMale   = true
	GenderFemale = false
)

type Customer struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;not null"`
	FirstName string     `gorm:"column:first_name;type:varchar(100);not null"`
	LastName  string     `gorm:"column:last_name;type:varchar(100);not null"`
	BirthDate time.Time  `gorm:"column:birth_date;type:date;not null"`
	Gender    bool       `gorm:"column:gender;type:boolean;not null"`
	Email     string     `gorm:"column:email;type:varchar(255);not null"`
	Address   *string    `gorm:"column:address;type:varchar(200);default:null"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp without time zone;not null"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp without time zone;not null"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp without time zone;default:null;index"`
}

func (c *Customer) FullName() string {
	return strings.TrimSpace(c.FirstName + " " + c.LastName)
}
