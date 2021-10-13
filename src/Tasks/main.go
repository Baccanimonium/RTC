package Tasks

import (
	"github.com/go-co-op/gocron"
	"time"
	"video-chat-app/src/Services"
	"video-chat-app/src/SocketHandlers"
)

type TaskManager struct {
	services  *Services.Services
	hub       *SocketHandlers.Hub
	scheduler *gocron.Scheduler
}

func NewTaskManager(services *Services.Services, hub *SocketHandlers.Hub) *TaskManager {
	return &TaskManager{services: services, hub: hub, scheduler: gocron.NewScheduler(time.UTC)}
}

func (tm *TaskManager) Run() {

}
