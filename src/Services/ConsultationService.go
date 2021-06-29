package Services

import "video-chat-app/src/Repos"

type ConsultationRepo struct {
	repo Repos.ConsultationRepo
}

func NewConsultationService(repo Repos.ConsultationRepo) *ConsultationRepo {
	return &ConsultationRepo{repo: repo}
}

func (s *ConsultationRepo) CreateConsultation(consultation Repos.Consultation) (int, error) {
	return s.repo.CreateConsultation(consultation)
}

func (s *ConsultationRepo) GetAllConsultation() ([]Repos.Consultation, error) {
	return s.repo.GetAllConsultation()
}

func (s *ConsultationRepo) GetConsultationById(id int) (Repos.Consultation, error) {
	return s.repo.GetConsultationById(id)
}
func (s *ConsultationRepo) UpdateConsultation(consultation Repos.Consultation, id int) (Repos.Consultation, error) {
	return s.repo.UpdateConsultation(consultation, id)
}

func (s *ConsultationRepo) DeleteConsultation(id int) error {
	return s.repo.DeleteConsultation(id)
}
