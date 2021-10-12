package Services

import (
	RTC "video-chat-app"
	"video-chat-app/src/Repos"
)

type ScheduleRepo struct {
	repo Repos.ScheduleRepo
	b    chan RTC.BroadcastingMessage
}

func NewScheduleService(repo Repos.ScheduleRepo, broadcast chan RTC.BroadcastingMessage) *ScheduleRepo {
	return &ScheduleRepo{repo: repo, b: broadcast}
}

func (s *ScheduleRepo) CreateSchedule(schedule Repos.Schedule) (Repos.Schedule, error) {
	newSchedule, err := s.repo.CreateSchedule(schedule)

	rawSchedule, convertError := RTC.ConvertToJson(newSchedule)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastCreateSchedule,
			Payload:     rawSchedule,
		}
	}

	return newSchedule, err
}

func (s *ScheduleRepo) UpdateSchedule(schedule Repos.Schedule, id int) (Repos.Schedule, error) {
	updatedSchedule, err := s.repo.UpdateSchedule(schedule, id)

	rawSchedule, convertError := RTC.ConvertToJson(updatedSchedule)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastUpdateSchedule,
			Payload:     rawSchedule,
		}
	}

	return updatedSchedule, err
}

func (s *ScheduleRepo) GetScheduleById(id int) (Repos.Schedule, error) {
	return s.repo.GetScheduleById(id)
}

func (s *ScheduleRepo) GetAllSchedule(idPatient int) ([]Repos.Schedule, error) {
	return s.repo.GetAllSchedule(idPatient)
}

func (s *ScheduleRepo) DeleteSchedule(id int) error {
	deletedSchedule, err := s.repo.DeleteSchedule(id)

	rawSchedule, convertError := RTC.ConvertToJson(deletedSchedule)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastDeleteSchedule,
			Payload:     rawSchedule,
		}
	}

	return err
}
