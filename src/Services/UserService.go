package Services

import "video-chat-app/src/Repos"

type UserRepo struct {
	repo Repos.UserRepo
}

func NewUserService(repo Repos.UserRepo) *UserRepo {
	return &UserRepo{repo: repo}
}

func (s *UserRepo) GetAllUser() ([]Repos.UserCreate, error) {
	return s.repo.GetAllUser()
}

func (s *UserRepo) GetUserById(id int) (Repos.UserCreate, error) {
	return s.repo.GetUserById(id)
}

func (s *UserRepo) UpdateUser(user Repos.UserCreate, id int) (Repos.UserCreate, error) {
	return s.repo.UpdateUser(user, id)
}

func (s *UserRepo) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
