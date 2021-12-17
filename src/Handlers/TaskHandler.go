package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src/Repos"
)

func (h Handler) CreateTask(c *gin.Context) {
	var input Repos.TaskCandidate

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskCandidatesList, err := h.services.TaskService.CreateTask(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, taskCandidatesList)
}

func (h Handler) GetAllTasks(c *gin.Context) {
	idDoctor, err := strconv.Atoi(c.Query("id_doctor"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid idDoctor param")
		return
	}
	idPatient, err := strconv.Atoi(c.Query("id_patient"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid idDoctor param")
		return
	}

	taskCandidatesList, err := h.services.TaskService.GetAllTasks(idDoctor, idPatient)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, taskCandidatesList)
}

func (h Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid task id param")
		return
	}

	if err := h.services.TaskCandidatesService.DeleteTaskCandidate(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
