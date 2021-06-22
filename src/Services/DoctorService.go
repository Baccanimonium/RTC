package Services

import "video-chat-app/src/Repos"

type DoctorRepo struct {
	repo Repos.DoctorRepo
}

func NewDoctorService(repo Repos.DoctorRepo) *DoctorRepo {
	return &DoctorRepo{repo: repo}
}

func (s *DoctorRepo) CreateDoctor(doctor Repos.Doctor) (int, error) {
	return s.repo.CreateDoctor(doctor)
}
