package models

import (
	"avitoTestTask/internal/models/helper"
	"github.com/google/uuid"
	"time"
)

type Bid struct {
	ID                 uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name               string           `gorm:"type:varchar(100)" json:"name"`
	Description        string           `gorm:"type:text" json:"description"`
	Status             string           `gorm:"type:varchar(20)" json:"status"`
	TenderID           uuid.UUID        `gorm:"type:uuid;index" json:"tenderId"`
	OrganizationID     uuid.UUID        `gorm:"type:uuid;index" json:"organizationId"` // ID организации, подавшей предложение
	CreatorUsername    string           `gorm:"type:varchar(50)" json:"creatorUsername"`
	RemainingApprovals int              `gorm:"default:3" json:"-"`  // Количество оставшихся согласий работников (не выводится в JSON)
	ApprovedUsers      helper.UUIDArray `gorm:"type:jsonb" json:"-"` // Список сотрудников, давших согласие (не выводится в JSON)
	CreatedAt          time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime" json:"-"`
	Version            uint             `gorm:"default:1" json:"version"`
	// ID организации, чей тендер
	OrganizationTenderID uuid.UUID `gorm:"type:uuid;index" json:"-"`
}

type BidVersion struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BidID       uuid.UUID `gorm:"type:uuid;index"`
	Version     uint      `gorm:"not null"`
	Name        string    `gorm:"type:varchar(100)"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"type:varchar(20)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type UpdateBidRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
