package Handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-chat-app/src/Models"
)

func (h Handler) getConsultation(c *gin.Context) {
	idConsultation, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid idConsultation param")
		return
	}

	consultation, err := h.services.ConsultationService.GetConsultationById(idConsultation)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, consultation)
}

func (h Handler) createConsultation(c *gin.Context) {
	var input Models.Consultation

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	consultation, err := h.services.ConsultationService.CreateConsultation(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, consultation)
}

func (h Handler) updateConsultation(c *gin.Context) {
	//idSchedule, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	newErrorResponse(c, http.StatusBadRequest, "invalid idSchedule param")
	//	return
	//}
	idConsultation, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid idConsultation param")
		return
	}

	var input Models.Consultation

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	consultation, err := h.services.ConsultationService.UpdateConsultation(input, idConsultation)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, consultation)

}

func (h Handler) listConsultation(c *gin.Context) {
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

	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid limit param")
		return
	}
	end, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid offset param")
		return
	}

	consultationList, err := h.services.ConsultationService.GetAllConsultation(Models.GetConsultationList{
		IdDoctor:  idDoctor,
		IdPatient: idPatient,
		Start:     start,
		End:       end,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, consultationList)
}

func (h Handler) deleteConsultation(c *gin.Context) {
	idConsultation, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid idConsultation param")
		return
	}

	if err := h.services.ConsultationService.DeleteConsultation(idConsultation); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": idConsultation})
}
