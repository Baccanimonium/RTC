package Tasks

import (
	"github.com/sirupsen/logrus"
	"time"
	RTC "video-chat-app"
)

func (tm *TaskManager) RunEventTask() error {
	_, err := tm.scheduler.Every(1).Minute().Tag("EventTasks").Do(tm.eventTask)
	return err
}

func (tm *TaskManager) eventTask() {
	events, error := tm.services.EventService.GetEventsByDate(time.Now().Format("02.01.2006 15:04"))
	if error != nil {
		logrus.Print("Failed event task ", error.Error())
	}

	for _, event := range events {
		rawEvent, convertError := RTC.ConvertToJson(event)
		if convertError == nil {
			broadcastingMessage := RTC.BroadcastingMessage{
				MessageType: RTC.BroadcastUpComingEvent,
				Payload:     rawEvent,
			}
			tm.hub.SendMessageToClient(broadcastingMessage, event.IdUser)
		}
	}
}
