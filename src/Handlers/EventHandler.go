package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src/Repos"
)

func (h Handler) createEvent(c *gin.Context) {
	var input Repos.Event

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.EventService.CreateEvent(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h Handler) listEvent(c *gin.Context) {
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

	scheduleList, err := h.services.EventService.GetAllEvents(Repos.GetAllEventsParams{
		IdDoctor:  idDoctor,
		IdPatient: idPatient,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, scheduleList)
}

func (h Handler) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	event, err := h.services.EventService.GetEventById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h Handler) UpdateEvent(c *gin.Context) {
	var input Repos.Event

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.services.EventService.UpdateEvent(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if _, err := h.services.EventService.DeleteEvent(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
