package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src/Models"
)

func (h Handler) createSchedule(c *gin.Context) {
	var input Models.DoctorSchedule

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
	scheduleList, err := h.services.ScheduleService.GetAllSchedule(Models.PostgresPagination{
		Limit: c.Param("limit"),
		Skip:  c.Param("offset"),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, scheduleList)
}

func (h Handler) getSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid doctor_id param")
		return
	}

	schedule, err := h.services.ScheduleService.GetScheduleByDoctorId(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) UpdateSchedule(c *gin.Context) {
	var input Models.DoctorSchedule

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.services.ScheduleService.UpdateSchedule(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h Handler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid doctor_id param")
		return
	}

	if err := h.services.ScheduleService.DeleteSchedule(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
