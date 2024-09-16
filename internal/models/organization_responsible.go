package models

import (
	"github.com/google/uuid"
)

type OrganizationResponsible struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid"`
	UserID         uuid.UUID `gorm:"type:uuid"`
}

func (OrganizationResponsible) TableName() string {
	return "organization_responsible"
}
