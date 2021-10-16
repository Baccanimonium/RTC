package Tasks

import (
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis"
	"time"
	"video-chat-app/src/Services"
	"video-chat-app/src/SocketHandlers"
)

type TaskManager struct {
	services  *Services.Services
	hub       *SocketHandlers.Hub
	rdb       *redis.Client
	scheduler *gocron.Scheduler
}

func NewTaskManager(services *Services.Services, hub *SocketHandlers.Hub, rdb *redis.Client) *TaskManager {
	return &TaskManager{services: services, hub: hub, scheduler: gocron.NewScheduler(time.UTC), rdb: rdb}
}

func (tm *TaskManager) Run() {

}
