package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	clients "insider_task/internal/client"
	"insider_task/internal/client/notification"
	"insider_task/internal/model"
	"insider_task/internal/repositories"
	"log"
	"time"
)

var (
	ticker     *time.Ticker
	tickerQuit = make(chan struct{})
	isRunning  = false
)

type MessageService struct {
	messageRepository *repositories.MessagesRepository
	clients           *clients.Clients
	redisClient       *redis.Client
}

func NewMessageService(messageRepository *repositories.MessagesRepository,
	clients *clients.Clients, redisClient *redis.Client) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		clients:           clients,
		redisClient:       redisClient,
	}
}

func (m *MessageService) StartMessageSender(ctx context.Context) error {
	if isRunning {
		log.Println("Message sender is already running")
		return nil
	}

	ticker = time.NewTicker(2 * time.Minute)
	isRunning = true
	log.Println("Started automatic message sender")

	go func() {
		for {
			select {
			case <-ticker.C:
				m.processUndeliveredMessages(ctx)
			case <-tickerQuit:
				ticker.Stop()
				isRunning = false
				log.Println("Stopped automatic message sender")
				return
			}
		}
	}()

	return nil
}

func (m *MessageService) processUndeliveredMessages(ctx context.Context) {
	messages, err := m.messageRepository.GetUndeliveredMessages()
	if err != nil {
		log.Printf("Error fetching undelivered messages: %v", err)
		return
	}

	for _, message := range messages {
		m.sendAndTrackMessage(ctx, message)
	}
}

func (m *MessageService) StopMessageSender() {
	if isRunning {
		tickerQuit <- struct{}{}
	}
}

func (m *MessageService) GetDeliveredMessages() ([]*model.MessageView, error) {
	dbMessages, err := m.messageRepository.GetDeliveredMessages()
	if err != nil {
		return nil, err
	}
	var messages = make([]*model.MessageView, len(dbMessages))
	for i, v := range dbMessages {
		messages[i] = &model.MessageView{
			ID:        v.ID,
			Phone:     v.Phone,
			Content:   v.Content,
			Sent:      v.IsSent,
			SentAt:    v.SentAt,
			MessageID: v.MessageID,
		}
	}
	return messages, nil
}

func (m *MessageService) SaveMessage(message *model.MessageDto) error {
	var messageEntity = &repositories.Message{
		Phone:   message.Phone,
		Content: message.Content,
		IsSent:  false,
	}
	return m.messageRepository.SaveMessage(messageEntity)
}

func (m *MessageService) sendAndTrackMessage(ctx context.Context, message *repositories.Message) {
	req := &notification.MessageRequest{
		To:      message.Phone,
		Content: truncateContent(message.Content, 160),
	}

	resp, err := m.clients.Notification.SendMessage(context.Background(), req)
	if err != nil {
		log.Printf("Failed to send message to %s: %v", message.Phone, err)
		return
	}

	if resp.Message != "Accepted" {
		log.Printf("Message to %s was not accepted", message.Phone)
		return
	}

	now := time.Now()
	message.SentAt = &now
	message.IsSent = true
	message.MessageID = &resp.MessageId

	if err := m.messageRepository.UpdateMessage(message); err != nil {
		log.Printf("Failed to update message ID %v: %v", message.ID, err)
	} else {
		log.Printf("Message ID %v marked as sent", message.ID)
	}

	m.cacheMessage(ctx, resp, now)
}

func (m *MessageService) cacheMessage(ctx context.Context, resp *notification.MessageResponse, now time.Time) {
	if resp == nil || resp.MessageId == "" {
		log.Println("[cacheMessage] Invalid or nil notification response")
		return
	}

	cacheKey := "message:" + resp.MessageId
	cacheValue := model.MessageCacheValue{
		MessageID: resp.MessageId,
		Timestamp: now.Format(time.RFC3339),
	}

	data, err := json.Marshal(cacheValue)
	if err != nil {
		log.Printf("[cacheMessage] Failed to marshal value for message %s: %v", resp.MessageId, err)
		return
	}

	if err := m.redisClient.Set(ctx, cacheKey, data, 2*time.Hour).Err(); err != nil {
		log.Printf("[cacheMessage] Failed to cache message %s to Redis: %v", resp.MessageId, err)
	}
}

func truncateContent(content string, limit int) string {
	if len(content) > limit {
		return content[:limit]
	}
	return content
}
