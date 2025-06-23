package router

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"insider_task/internal/configs"
	"insider_task/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	config         *configs.Router
	engine         *gin.Engine
	messageHandler *handlers.MessageHandler
}

func NewRouter(config *configs.Router, messageHandler *handlers.MessageHandler) *Router {
	return &Router{
		config:         config,
		engine:         gin.Default(),
		messageHandler: messageHandler,
	}
}

func (r *Router) InitAndRun() error {
	root := r.engine.Group("/api/v1/messages")
	r.registerMessagesRoutes(root)
	r.registerSwaggerHandler(root)

	return r.engine.Run(":" + r.config.Port)
}

func (r *Router) registerMessagesRoutes(root *gin.RouterGroup) {
	messages := root.Group("")
	{
		messages.POST("/start", r.messageHandler.StartAutoSender)
		messages.POST("/stop", r.messageHandler.StopAutoSender)
		messages.POST("", r.messageHandler.CreateMessage)
		messages.GET("/sent", r.messageHandler.GetDeliveredMessages)
	}
}

func (r *Router) registerSwaggerHandler(v1 *gin.RouterGroup) {
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
