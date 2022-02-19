package Services

import (
	RTC "video-chat-app"
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type ScheduleRepo struct {
	repo Repos.ScheduleRepo
	b    chan RTC.BroadcastingMessage
}

func NewScheduleService(repo Repos.ScheduleRepo, broadcast chan RTC.BroadcastingMessage) *ScheduleRepo {
	return &ScheduleRepo{repo: repo, b: broadcast}
}

func (s *ScheduleRepo) CreateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error) {
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

func (s *ScheduleRepo) UpdateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error) {
	updatedSchedule, err := s.repo.UpdateSchedule(schedule)

	rawSchedule, convertError := RTC.ConvertToJson(updatedSchedule)

	if err == nil && convertError == nil {
		s.b <- RTC.BroadcastingMessage{
			MessageType: RTC.BroadcastUpdateSchedule,
			Payload:     rawSchedule,
		}
	}

	return updatedSchedule, err
}

func (s *ScheduleRepo) GetScheduleByDoctorId(id int) (Models.DoctorSchedule, error) {
	return s.repo.GetScheduleByDoctorId(id)
}

func (s *ScheduleRepo) GetAllSchedule(params Models.PostgresPagination) ([]Models.DoctorSchedule, error) {
	return s.repo.GetAllSchedule(params)
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
