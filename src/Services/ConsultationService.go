package Services

import "video-chat-app/src/Repos"

type ConsultationRepo struct {
	repo Repos.ConsultationRepo
}

func NewConsultationService(repo Repos.ConsultationRepo) *ConsultationRepo {
	return &ConsultationRepo{repo: repo}
}

func (s *ConsultationRepo) CreateConsultation(idSchedule int, consultation Repos.Consultation) (Repos.Consultation, error) {
	return s.repo.CreateConsultation(idSchedule, consultation)
}

func (s *ConsultationRepo) GetAllConsultation(idSchedule int) ([]Repos.Consultation, error) {
	return s.repo.GetAllConsultation(idSchedule)
}

func (s *ConsultationRepo) GetConsultationById(idSchedule, idConsultation int) (Repos.Consultation, error) {
	return s.repo.GetConsultationById(idSchedule, idConsultation)
}
func (s *ConsultationRepo) UpdateConsultation(consultation Repos.Consultation, idSchedule, idConsultation int) (Repos.Consultation, error) {
	return s.repo.UpdateConsultation(consultation, idSchedule, idConsultation)
}

func (s *ConsultationRepo) DeleteConsultation(idSchedule, idConsultation int) error {
	return s.repo.DeleteConsultation(idSchedule, idConsultation)
}
