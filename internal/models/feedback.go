package models

import (
	"github.com/google/uuid"
	"time"
)

type Feedback struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BidID uuid.UUID `gorm:"type:uuid;index"`
	// ID того, на кого оставляют отзыв (!)
	CreatorUsername string    `gorm:"type:varchar(50)"`
	FeedbackText    string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
