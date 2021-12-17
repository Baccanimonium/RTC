package Services

import "video-chat-app/src/Repos"

type TaskCandidates struct {
	repo Repos.TaskCandidatesRepo
}

func NewTaskCandidatesService(repo Repos.TaskCandidatesRepo) *TaskCandidates {
	return &TaskCandidates{repo: repo}
}

func (s *TaskCandidates) CreateTaskCandidates(unConfirmedEvents []Repos.Event, currentDate, currentTime string) error {
	return s.repo.CreateTaskCandidates(unConfirmedEvents, currentDate, currentTime)
}

func (s *TaskCandidates) DeleteTaskCandidate(taskId int) error {
	return s.repo.DeleteTaskCandidate(taskId)
}

func (s *TaskCandidates) ExtractTaskCandidates(taskTime string) ([]Repos.TaskCandidate, error) {
	return s.repo.ExtractTaskCandidates(taskTime)
}

func (s *TaskCandidates) GetTaskCandidatesByPatient(patientId int) ([]Repos.TaskCandidate, error) {
	return s.repo.GetTaskCandidatesByPatient(patientId)
}
