package Handlers

import (
	"github.com/gin-gonic/gin"
	"video-chat-app/src/Services"
)

type Handler struct {
	services *Services.Services
}

func NewHandler(services *Services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		doctors := api.Group("/doctors")
		{
			doctors.POST("/", h.createDoctor)
			doctors.GET("/", h.listDoctor)
			doctors.GET("/:id", h.getDoctor)
			doctors.PUT("/:id", h.updateDoctor)
			doctors.DELETE("/:id", h.deleteDoctor)
		}
		patient := api.Group("/patient")
		{
			patient.POST("/", h.createPatient)
			patient.GET("/", h.listPatient)
			patient.GET("/:id", h.getPatient)
			patient.PUT("/:id", h.UpdatePatient)
			patient.DELETE("/:id", h.DeletePatient)
		}
		schedule := api.Group("/schedule")
		{
			schedule.POST("/", h.createSchedule)
			schedule.GET("/", h.listSchedule)
			schedule.GET("/:id", h.getSchedule)
			schedule.PUT("/:id", h.UpdateSchedule)
			schedule.DELETE("/:id", h.DeleteSchedule)
			event := schedule.Group(":id/event")
			{
				event.POST("/", h.createEvent)
				event.GET("/", h.listEvent)
			}
		}

		event := schedule.Group("/event")
		{
			event.GET("/:id", h.getEvent)
			event.PUT("/:id", h.UpdateEvent)
			event.DELETE("/:id", h.DeleteEvent)
		}

	}

	return router
}
