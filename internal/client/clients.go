package clients

import (
	"context"
	"insider_task/internal/client/notification"
	"insider_task/internal/configs"
)

type NotificationClient interface {
	SendMessage(context.Context, *notification.MessageRequest) (*notification.MessageResponse, error)
}

var (
	_ NotificationClient = (*notification.Client)(nil)
)

type Clients struct {
	Notification NotificationClient
}

func NewClients(configs *configs.Configs) *Clients {
	return &Clients{
		Notification: notification.NewClient(configs.NotificationClient),
	}
}
