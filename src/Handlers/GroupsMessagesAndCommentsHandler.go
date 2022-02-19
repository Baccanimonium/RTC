package Handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"video-chat-app/src/Models"
)

func (h Handler) getGroupMessage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id_message"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	doctor, err := h.services.GroupMessagesService.GetGroupMessage(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h Handler) listGroupMessage(c *gin.Context) {
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
	doctorList, err := h.services.GroupMessagesService.GetGroupMessages(Models.GetGroupMessages{
		GroupId: c.Param("groupId"),
		Pagination: Models.MongoPagination{
			Limit: &limit,
			Skip:  &skip,
		},
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, doctorList)
}

func (h Handler) createGroupMessage(c *gin.Context) {
	var input Models.GroupMessage

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupMessagesService.CreateGroupMessage(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)
}

func (h Handler) updateGroupMessage(c *gin.Context) {
	var input Models.GroupMessage

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupMessagesService.UpdateGroupMessage(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)

}

func (h Handler) deleteGroupMessage(c *gin.Context) {
	if err := h.services.GroupMessagesService.DeleteGroupMessage(c.Param("message_id")); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) getGroupMessageComment(c *gin.Context) {

	doctor, err := h.services.GroupMessagesService.GetGroupMessageComment(c.Param("id_message_comment"))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h Handler) listGroupMessageComment(c *gin.Context) {
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
	doctorList, err := h.services.GroupMessagesService.GetGroupMessagesComment(Models.GetGroupMessagesComments{
		MessageId: c.Param("id_message"),
		Pagination: Models.MongoPagination{
			Limit: &limit,
			Skip:  &skip,
		},
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, doctorList)
}

func (h Handler) createGroupMessageComment(c *gin.Context) {
	var input Models.GroupMessageComment

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupMessagesService.CreateGroupMessageComment(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)
}

func (h Handler) updateGroupMessageComment(c *gin.Context) {
	var input Models.GroupMessageComment

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newGroup, err := h.services.GroupMessagesService.UpdateGroupMessageComment(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGroup)

}

func (h Handler) deleteGroupMessageComment(c *gin.Context) {
	if _, err := h.services.GroupMessagesService.DeleteGroupMessageComment(c.Param("id")); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
