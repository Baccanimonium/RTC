package Handlers

import (
	"github.com/gin-gonic/gin"
	"video-chat-app/src/Services"
)

type Handler struct {
	services *Services.Services
}

func NewHandler(services *Services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.singUp)
		auth.POST("/sing-in", h.singIn)
	}

	return router
}
