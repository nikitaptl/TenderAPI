package models

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Organization) TableName() string {
	return "organization"
}
