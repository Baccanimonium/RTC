package Handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"video-chat-app/src/Services"
	"video-chat-app/src/SocketHandlers"
)

type Handler struct {
	services *Services.Services
}

func NewHandler(services *Services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter(socketFactory *SocketHandlers.SocketClientFactory) *gin.Engine {
	router := gin.New()
	router.Use(cors.Default())
	router.GET("/websocket", socketFactory.OnNewSocketClient)
	//router.GET("/websocket", socketFactory.OnNewSocketClient)

	router.StaticFS("/file", http.Dir("public"))

	//router.NoRoute(func(c *gin.Context) {
	//   c.File("./public/index.html")
	//})

	router.POST("/upload", h.handleUploadFile)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.StaticFS("/file", http.Dir("public"))

		api.POST("/upload", h.handleUploadFile)

		channels := api.Group("/channels")
		{
			channels.GET("/", h.GetAllChannelsBelongsToUser)
		}
		messages := api.Group("/messages")
		{
			messages.GET("/:channelId", h.GetAllMessagesBelongsToChannel)
		}
		users := api.Group("/users")
		{
			users.POST("/", h.signUp)
			users.GET("/", h.listUser)
			users.GET("/profile", h.getUserProfile)
			users.GET("/:id", h.getUser)
			users.PUT("/:id", h.updateUser)
			users.DELETE("/:id", h.deleteUser)
		}
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
