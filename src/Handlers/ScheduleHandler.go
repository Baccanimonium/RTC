package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src/Repos"
)

func (h Handler) createSchedule(c *gin.Context) {
	var input Repos.Schedule

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.services.ScheduleService.CreateSchedule(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) listSchedule(c *gin.Context) {
	scheduleList, err := h.services.ScheduleService.GetAllSchedule()

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, scheduleList)
}

func (h Handler) getSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	schedule, err := h.services.ScheduleService.GetScheduleById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input Repos.Schedule

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.services.ScheduleService.UpdateSchedule(input, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.ScheduleService.DeleteSchedule(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}
