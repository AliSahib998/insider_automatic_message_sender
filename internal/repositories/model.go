package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Phone     string         `gorm:"column:phone;type:varchar(32);not null"`
	Content   string         `gorm:"column:content;type:text;not null"`
	IsSent    bool           `gorm:"column:is_sent;type:boolean;default:false;not null"`
	SentAt    *time.Time     `gorm:"column:sent_at"`
	MessageID *string        `gorm:"column:message_id;type:varchar(255)"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
