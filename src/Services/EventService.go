package Services

import "video-chat-app/src/Repos"

type EventRepo struct {
	repo Repos.EventRepo
}

func NewEventService(repo Repos.EventRepo) *EventRepo {
	return &EventRepo{repo: repo}
}

func (s *EventRepo) CreateEvent(idSchedule int, event Repos.Event) (int, error) {
	return s.repo.CreateEvent(idSchedule, event)
}

func (s *EventRepo) UpdateEvent(event Repos.Event, id int) (Repos.Event, error) {
	return s.repo.UpdateEvent(event, id)
}

func (s *EventRepo) GetEventById(id int) (Repos.Event, error) {
	return s.repo.GetEventById(id)
}

func (s *EventRepo) GetAllEvents(idSchedule int) ([]Repos.Event, error) {
	return s.repo.GetAllEvents(idSchedule)
}

func (s *EventRepo) DeleteEvent(id int) error {
	return s.repo.DeleteEvent(id)
}
