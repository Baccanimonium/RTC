package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"video-chat-app/src"
)

func (h Handler) GetAllChannelsBelongsToUser(c *gin.Context) {
	userId, userCtxError := c.Get(src.UserContext)

	if !userCtxError {
		c.JSON(http.StatusBadRequest, "user does not exist")
		return
	}

	channels, err := h.services.GetAllChannelsBelongsToUser(userId.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, channels)
}
