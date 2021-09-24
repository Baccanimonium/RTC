package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"video-chat-app/src/Services"
)

const (
	authorizationHeader = "Authorization"
	UserContext         = "IdUser"
)

func (h Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	IdUser, err := Services.ParseToken(headerParts[1])

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(UserContext, IdUser)
}
