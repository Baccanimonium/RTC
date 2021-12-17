package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src"
)

func (h Handler) getTaskCandidatesByPatientId(c *gin.Context) {
	userId, _ := c.Get(src.UserContext)

	taskCandidatesList, err := h.services.TaskCandidatesService.GetTaskCandidatesByPatient(userId.(int))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, taskCandidatesList)
}

func (h Handler) deleteTaskCandidate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid taskCandidate id param")
		return
	}

	if err := h.services.TaskCandidatesService.DeleteTaskCandidate(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
