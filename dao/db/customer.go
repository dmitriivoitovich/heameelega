package db

import (
	"time"

	"github.com/google/uuid"
)

const (
	GenderMale   = true
	GenderFemale = false
)

type Customer struct {
	BaseModel
	FirstName string    `gorm:"column:first_name;type:varchar(100);not null"`
	LastName  string    `gorm:"column:last_name;type:varchar(100);not null"`
	BirthDate time.Time `gorm:"column:birth_date;type:date;not null"`
	Gender    bool      `gorm:"column:gender;type:boolean;not null"`
	Email     string    `gorm:"column:email;type:varchar(255);not null;index:idx_customers_email_user_id_unique,unique"`
	Address   *string   `gorm:"column:address;type:varchar(200);default:null"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;index:idx_customers_email_user_id_unique,unique;index"`
}
