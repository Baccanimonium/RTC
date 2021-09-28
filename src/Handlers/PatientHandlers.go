package Handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"video-chat-app/src"
	"video-chat-app/src/Repos"
)

func (h Handler) createPatient(c *gin.Context) {
	var input Repos.Patient
	userId, _ := c.Get(src.UserContext)

	if err := c.BindJSON(&input); err != nil {
		logrus.Print("fail to get patient", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// adding current userId from token
	input.IdCurrentDoctor = userId.(int)

	id, err := h.services.PatientService.CreatePatient(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h Handler) getPatient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	patient, err := h.services.PatientService.GetPatientById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h Handler) UpdatePatient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input Repos.Patient

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	patient, err := h.services.PatientService.UpdatePatient(input, id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h Handler) DeletePatient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.PatientService.DeletePatient(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (h Handler) listPatient(c *gin.Context) {
	userId, _ := c.Get(src.UserContext)
	patientList, err := h.services.PatientService.GetAllPatient(userId.(int))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, patientList)
}
