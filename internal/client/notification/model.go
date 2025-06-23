package notification

type MessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type MessageResponse struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId"`
}
