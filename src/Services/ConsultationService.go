package Services

import (
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type ConsultationRepo struct {
	repo Repos.ConsultationRepo
}

func NewConsultationService(repo Repos.ConsultationRepo) *ConsultationRepo {
	return &ConsultationRepo{repo: repo}
}

func (s *ConsultationRepo) CreateConsultation(consultation Models.Consultation) (Models.Consultation, error) {
	return s.repo.CreateConsultation(consultation)
}

func (s *ConsultationRepo) GetAllConsultation(params Models.GetConsultationList) ([]Models.Consultation, error) {
	return s.repo.GetAllConsultation(params)
}

func (s *ConsultationRepo) GetConsultationById(idConsultation int) (Models.Consultation, error) {
	return s.repo.GetConsultationById(idConsultation)
}

func (s *ConsultationRepo) UpdateConsultation(consultation Models.Consultation, id int) (Models.Consultation, error) {
	return s.repo.UpdateConsultation(consultation, id)
}

func (s *ConsultationRepo) SetDoctorJoinTime(id int) error {
	return s.repo.SetDoctorJoinTime(id)
}

func (s *ConsultationRepo) DeleteConsultation(idConsultation int) error {
	return s.repo.DeleteConsultation(idConsultation)
}

func (s *ConsultationRepo) CreateConsultationNotes(notes Models.Notes) (Models.Notes, error) {
	return s.repo.CreateConsultationNotes(notes)
}
func (s *ConsultationRepo) UpdateConsultationNotes(notes Models.Notes) (Models.Notes, error) {
	return s.repo.UpdateConsultationNotes(notes)
}

func (s *ConsultationRepo) DeleteConsultationNotes(idNotes int) error {
	return s.repo.DeleteConsultationNotes(idNotes)
}
