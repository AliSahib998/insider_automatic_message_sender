package model

import (
	"github.com/google/uuid"
	"time"
)

type MessageView struct {
	ID        uuid.UUID  `json:"id"`
	Phone     string     `json:"phone"`
	Content   string     `json:"content"`
	Sent      bool       `json:"sent"`
	SentAt    *time.Time `json:"sentAt"`
	MessageID *string    `json:"messageId"`
}

type MessageDto struct {
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

type MessageCacheValue struct {
	MessageID string `json:"messageId"`
	Timestamp string `json:"timestamp"`
}
