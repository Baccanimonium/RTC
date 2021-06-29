package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h Handler) getConsultation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	consultation, err := h.services.DoctorService.GetDoctorById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, consultation)
}

func (h Handler) createConsultation(c *gin.Context) {

}

func (h Handler) updateConsultation(c *gin.Context) {

}

func (h Handler) lostConsultation(c *gin.Context) {

}

func (h Handler) deleteConsultation(c *gin.Context) {

}
