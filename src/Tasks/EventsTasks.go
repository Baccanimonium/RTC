package Tasks

import (
	"github.com/sirupsen/logrus"
	"time"
	RTC "video-chat-app"
	"video-chat-app/src/Repos"
)

func (tm *TaskManager) RunEventTask() error {
	_, err := tm.scheduler.Every(1).Minute().Tag("EventTasks").Do(tm.eventTask)
	return err
}

func (tm *TaskManager) eventTask() {
	currentDate := time.Now()
	currentTime := currentDate.Format("15:04")
	currentDateString := currentDate.Format("02.01.2006")
	events, err := tm.services.EventService.GetAllEvents(Repos.GetAllEventsParams{Date: currentDate})
	if err != nil {
		logrus.Print("Failed event task ", err.Error())
	}

	unConfirmedEvents := make([]Repos.Event, 0)

	for _, event := range events {
		rawEvent, convertError := RTC.ConvertToJson(event)
		if convertError == nil {
			broadcastingMessage := RTC.BroadcastingMessage{
				MessageType: RTC.BroadcastUpComingEvent,
				Payload:     rawEvent,
			}
			tm.hub.SendMessageToClient(broadcastingMessage, event.IdPatient)
			if event.NotifyDoctor {
				tm.hub.SendMessageToClient(broadcastingMessage, event.IdDoctor)
			}
		}
		if event.RequiresConfirmation {
			unConfirmedEvents = append(unConfirmedEvents, event)
		}
	}

	if len(unConfirmedEvents) > 0 {
		err := tm.services.TaskCandidatesService.CreateTaskCandidates(unConfirmedEvents, currentDateString, currentTime)

		if err != nil {
			logrus.Print("Failed To create TASK CANDIDATES LIST AT ", currentTime, " ERROR ", err.Error())
		}
	}
}
