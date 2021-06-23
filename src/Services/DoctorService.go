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

func (s *DoctorRepo) GetAllDoctor() ([]Repos.Participant, error) {
	return s.repo.GetAllDoctor()
}

func (s *DoctorRepo) GetDoctorById(id int) (Repos.Participant, error) {
	return s.repo.GetDoctorById(id)
}

func (s *DoctorRepo) UpdateDoctor(doctor Repos.Doctor, id int) (Repos.Doctor, error) {
	return s.repo.UpdateDoctor(doctor, id)
}

func (s *DoctorRepo) DeleteDoctor(id int) error {
	return s.repo.DeleteDoctor(id)
}
