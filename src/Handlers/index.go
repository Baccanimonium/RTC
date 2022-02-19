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

	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

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

		taskCandidates := api.Group("/task_candidates")
		{
			taskCandidates.GET("/", h.getTaskCandidatesByPatientId)
			taskCandidates.DELETE("/:id", h.deleteTaskCandidate)
		}

		task := api.Group("/task")
		{
			task.GET("/", h.GetAllTasks)
			task.POST("/", h.CreateTask)
			task.DELETE("/:id", h.DeleteTask)
		}

		event := api.Group("/event")
		{
			event.POST("/", h.createEvent)
			event.GET("/", h.listEvent)
			event.GET("/:id", h.getEvent)
			event.PUT("/:id", h.UpdateEvent)
			event.DELETE("/:id", h.DeleteEvent)
		}

		consultation := api.Group("/consultation")
		{
			consultation.POST("/", h.createConsultation)
			consultation.GET("/", h.listConsultation)
			consultation.GET("/:id", h.getConsultation)
			consultation.PUT("/:id", h.updateConsultation)
			consultation.DELETE("/:id", h.deleteConsultation)
		}

		schedule := api.Group("/schedule")
		{
			schedule.POST("/", h.createSchedule)
			schedule.GET("/", h.listSchedule)
			schedule.GET("/:doctor_id", h.getSchedule)
			schedule.PUT("/:doctor_id", h.UpdateSchedule)
			schedule.DELETE("/:doctor_id", h.DeleteSchedule)
		}

		group := api.Group("/group")
		{
			group.POST("/", h.createGroup)
			group.GET("/", h.getGroups)
			group.GET("/:id", h.getGroup)
			group.PUT("/:id", h.updateGroup)
			group.DELETE("/:id", h.deleteGroup)
			group.POST("/subscribe", h.subscribeToGroup)
			group.POST("/un_subscribe", h.unSubscribeToGroup)
			group.POST("/pin_group_message", h.pinGroupMessage)

			groupMessages := group.Group(":id/messages")
			{
				groupMessages.POST("/", h.createGroupMessage)
				groupMessages.GET("/", h.listGroupMessage)
				//groupMessages.GET("/:message_id", h.getConsultation)
				groupMessages.PUT("/:message_id", h.updateGroupMessage)
				groupMessages.DELETE("/:message_id", h.deleteGroupMessage)
			}
			groupMessagesComments := groupMessages.Group(":message_id/comments")
			{
				groupMessagesComments.POST("/", h.createGroupMessageComment)
				groupMessagesComments.GET("/", h.listGroupMessageComment)
				groupMessagesComments.GET("/:comment_id", h.getGroupMessageComment)
				groupMessagesComments.PUT("/:comment_id", h.updateGroupMessageComment)
				groupMessagesComments.DELETE("/:comment_id", h.deleteGroupMessageComment)
			}
		}
	}

	return router
}
