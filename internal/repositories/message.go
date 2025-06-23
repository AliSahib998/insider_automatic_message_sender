package repositories

import (
	"gorm.io/gorm"
)

type MessagesRepository struct {
	db *gorm.DB
}

func NewMessagesRepository(db *gorm.DB) *MessagesRepository {
	return &MessagesRepository{db: db}
}

func (r *MessagesRepository) GetDeliveredMessages() ([]*Message, error) {
	var messages []*Message
	result := r.db.Where("is_sent = ?", true).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (r *MessagesRepository) GetUndeliveredMessages() ([]*Message, error) {
	var messages []*Message
	result := r.db.Where("is_sent = ?", false).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (r *MessagesRepository) SaveMessage(msg *Message) error {
	result := r.db.Create(msg)
	return result.Error
}

func (r *MessagesRepository) UpdateMessage(message *Message) error {
	return r.db.Model(&Message{}).
		Where("id = ?", message.ID).
		Updates(map[string]interface{}{
			"is_sent":    message.IsSent,
			"sent_at":    message.SentAt,
			"message_id": message.MessageID,
		}).Error
}
