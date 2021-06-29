package Services

import "video-chat-app/src/Repos"

type ScheduleRepo struct {
	repo Repos.ScheduleRepo
}

func NewScheduleService(repo Repos.ScheduleRepo) *ScheduleRepo {
	return &ScheduleRepo{repo: repo}
}

func (s *ScheduleRepo) CreateSchedule(schedule Repos.Schedule) (int, error) {
	return s.repo.CreateSchedule(schedule)
}

func (s *ScheduleRepo) GetAllSchedule() ([]Repos.Schedule, error) {
	return s.repo.GetAllSchedule()
}

func (s *ScheduleRepo) GetScheduleById(id int) (Repos.Schedule, error) {
	return s.repo.GetScheduleById(id)
}
func (s *ScheduleRepo) UpdateSchedule(schedule Repos.Schedule, id int) (Repos.Schedule, error) {
	return s.repo.UpdateSchedule(schedule, id)
}

func (s *ScheduleRepo) DeleteSchedule(id int) error {
	return s.repo.DeleteSchedule(id)
}
