package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"video-chat-app/src/Repos"
)

func (h Handler) getDoctor(c *gin.Context) {

}

func (h Handler) createDoctor(c *gin.Context) {
	var input Repos.Doctor

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.DoctorService.CreateDoctor(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h Handler) updateDoctor(c *gin.Context) {

}

func (h Handler) deleteDoctor(c *gin.Context) {

}

func (h Handler) listDoctor(c *gin.Context) {
	doctorList, err := h.services.DoctorService.GetAllDoctor()

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, doctorList)
}
