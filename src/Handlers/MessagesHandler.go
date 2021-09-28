package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"video-chat-app/src"
)

func (h Handler) GetAllMessagesBelongsToChannel(c *gin.Context) {
	userId, userCtxError := c.Get(src.UserContext)
	channelId := c.Param("channelId")

	if !userCtxError {
		c.JSON(http.StatusBadRequest, "user does not exist")
		return
	}

	channels, err := h.services.MessagesService.GetMessages(channelId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, channels)
}
