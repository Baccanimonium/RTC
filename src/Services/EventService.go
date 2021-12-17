package Services

import "video-chat-app/src/Repos"

type EventRepo struct {
	repo Repos.EventRepo
}

func NewEventService(repo Repos.EventRepo) *EventRepo {
	return &EventRepo{repo: repo}
}

func (s *EventRepo) CreateEvent(event Repos.Event) (Repos.Event, error) {
	return s.repo.CreateEvent(event)
}

func (s *EventRepo) UpdateEvent(event Repos.Event) (Repos.Event, error) {
	return s.repo.UpdateEvent(event)
}

func (s *EventRepo) GetEventById(id int) (Repos.Event, error) {
	return s.repo.GetEventById(id)
}

func (s *EventRepo) GetAllEvents(request Repos.GetAllEventsParams) ([]Repos.Event, error) {
	return s.repo.GetAllEvents(request)
}

func (s *EventRepo) DeleteEvent(id int) (Repos.Event, error) {
	return s.repo.DeleteEvent(id)
}
