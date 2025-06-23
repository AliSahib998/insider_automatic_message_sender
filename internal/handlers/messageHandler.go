package handlers

import (
	"insider_task/internal/model"
	"insider_task/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MessageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		MessageService: messageService,
	}
}

// StartAutoSender starts the automatic message sending ticker
// @Summary      Start automatic message sending
// @Description  Starts a background process that sends unsent messages every 2 minutes
// @Tags         Message
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "status message"
// @Failure      500 {object} map[string]string "error message"
// @Router       /start [post]
func (r *MessageHandler) StartAutoSender(c *gin.Context) {
	err := r.MessageService.StartMessageSender(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "started"})
}

// StopAutoSender stops the automatic message sending ticker
// @Summary      Stop automatic message sending
// @Description  Stops the background process that sends unsent messages every 2 minutes
// @Tags         Message
// @Accept       json
// @Produce      json
// @Success      204 "No Content"
// @Router       /stop [post]
func (r *MessageHandler) StopAutoSender(c *gin.Context) {
	r.MessageService.StopMessageSender()
	c.Status(http.StatusNoContent)
}

// GetDeliveredMessages returns the list of messages that have been sent
// @Summary      List delivered (sent) messages
// @Description  Retrieves all messages from the database that have been marked as sent
// @Tags         Message
// @Accept       json
// @Produce      json
// @Success      200 {array} model.MessageView "List of delivered messages"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /sent [get]
func (r *MessageHandler) GetDeliveredMessages(c *gin.Context) {
	response, err := r.MessageService.GetDeliveredMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// CreateMessage creates a new message to be sent later
// @Summary      Create message
// @Description  Save a new message with phone number and content. It will be sent automatically later.
// @Tags         Message
// @Accept       json
// @Produce      json
// @Param        message  body  model.MessageDto  true  "Message data"
// @Success      201 {object} map[string]string "Created successfully"
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       / [post]
func (r *MessageHandler) CreateMessage(c *gin.Context) {
	var dto model.MessageDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := r.MessageService.SaveMessage(&dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}
