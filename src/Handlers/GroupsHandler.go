package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"video-chat-app/src/Models"
)

func (h Handler) getGroup(c *gin.Context) {

	group, err := h.services.GroupsService.GetGroup(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h Handler) getGroups(c *gin.Context) {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid limit param")
		return
	}
	skip, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid offset param")
		return
	}

	groupsList, err := h.services.GroupsService.GetGroups(Models.GetGroupFilterParams{
		Tags: strings.Split(c.Param("tags"), ","),
		Pagination: Models.MongoPagination{
			Limit: &limit,
			Skip:  &skip,
		},
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, groupsList)
}

func (h Handler) createGroup(c *gin.Context) {
	var input Models.Group

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupsService.CreateGroup(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)
}

func (h Handler) updateGroup(c *gin.Context) {
	var input Models.Group

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupsService.UpdateGroup(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)

}

func (h Handler) deleteGroup(c *gin.Context) {
	if _, err := h.services.GroupsService.DeleteGroup(c.Param("id")); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) subscribeToGroup(c *gin.Context) {

	var input Models.GroupSubscription

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subscriptions, err := h.services.GroupsService.SubscribeToGroup(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (h Handler) unSubscribeToGroup(c *gin.Context) {
	var input Models.GroupSubscription

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	subscriptions, err := h.services.GroupsService.UnSubscribeToGroup(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (h Handler) pinGroupMessage(c *gin.Context) {
	var input Models.GroupPinMessage

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.GroupsService.PinGroupMessage(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
