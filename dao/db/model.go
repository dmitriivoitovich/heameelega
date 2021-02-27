package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"column:id;primaryKey;type:uuid;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp without time zone;not null;index"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp without time zone;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp without time zone;default:null;index"`
}

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

	if m.DeletedAt.Valid {
		tx.Statement.SetColumn("DeletedAt", m.DeletedAt)
	}

	return nil
}
