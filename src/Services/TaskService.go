package Services

import "video-chat-app/src/Repos"

type Tasks struct {
	repo Repos.TaskRepo
}

func NewTaskService(repo Repos.TaskRepo) *Tasks {
	return &Tasks{repo: repo}
}

func (s *Tasks) CreateTask(task Repos.TaskCandidate) (Repos.Task, error) {
	return s.repo.CreateTask(task)
}

func (s *Tasks) GetAllTasks(idDoctor int, idPatient int) ([]Repos.Task, error) {
	return s.repo.GetAllTasks(idDoctor, idPatient)
}

func (s *Tasks) DeleteTask(idTask int) (Repos.Task, error) {
	return s.repo.DeleteTask(idTask)
}
