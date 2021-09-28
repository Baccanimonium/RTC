package Services

import "video-chat-app/src/Repos"

type PatientRepo struct {
	repo Repos.PatientRepo
}

func NewPatientService(repo Repos.PatientRepo) *PatientRepo {
	return &PatientRepo{repo: repo}
}

func (s *PatientRepo) CreatePatient(patient Repos.Patient) (int, error) {
	return s.repo.CreatePatient(patient)
}

func (s *PatientRepo) GetAllPatient(userId int) ([]Repos.Participant, error) {
	return s.repo.GetAllPatient(userId)
}

func (s *PatientRepo) GetPatientById(id int) (Repos.Participant, error) {
	return s.repo.GetPatientById(id)
}

func (s *PatientRepo) UpdatePatient(patient Repos.Patient, id int) (Repos.Patient, error) {
	return s.repo.UpdatePatient(patient, id)
}

func (s *PatientRepo) DeletePatient(id int) error {
	return s.repo.DeletePatient(id)
}
