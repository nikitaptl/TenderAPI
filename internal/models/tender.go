package models

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"type:varchar(100)" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	Status          string    `gorm:"type:varchar(20)" json:"status"`
	ServiceType     string    `gorm:"type:varchar(50)" json:"serviceType"`
	Version         uint      `gorm:"default:1" json:"version"`
	OrganizationID  uuid.UUID `gorm:"type:uuid;index" json: "organizationId"`
	CreatorUsername string    `gorm:"type:varchar(50)" json:"creatorUsername"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type TenderVersion struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenderID    uuid.UUID `gorm:"type:uuid"`
	Version     uint      `gorm:"not null"`
	Name        string    `gorm:"type:varchar(100)"`
	Description string    `gorm:"type:text"`
	ServiceType string    `gorm:"type:varchar(50)"`
	Status      string    `gorm:"type:varchar(20)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

/* Для строгого конроля полей, которые пользователь может обновить
(Чтобы он не мог, например, обновить ID) создаём отдельную структуру */

type UpdateTenderRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ServiceType *string `json:"serviceType,omitempty"`
}
