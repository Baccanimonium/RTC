package Services

import "video-chat-app/src/Repos"

type ConsultationRepo struct {
	repo Repos.ConsultationRepo
}

func NewConsultationService(repo Repos.ConsultationRepo) *ConsultationRepo {
	return &ConsultationRepo{repo: repo}
}

func (s *ConsultationRepo) CreateConsultation(consultation Repos.Consultation) (Repos.Consultation, error) {
	return s.repo.CreateConsultation(consultation)
}

func (s *ConsultationRepo) GetAllConsultation(idDoctor int, idPatient int) ([]Repos.Consultation, error) {
	return s.repo.GetAllConsultation(idDoctor, idPatient)
}

func (s *ConsultationRepo) GetConsultationById(idConsultation int) (Repos.Consultation, error) {
	return s.repo.GetConsultationById(idConsultation)
}

func (s *ConsultationRepo) UpdateConsultation(consultation Repos.Consultation, id int) (Repos.Consultation, error) {
	return s.repo.UpdateConsultation(consultation, id)
}

func (s *ConsultationRepo) SetDoctorJoinTime(id int) error {
	return s.repo.SetDoctorJoinTime(id)
}

func (s *ConsultationRepo) DeleteConsultation(idConsultation int) error {
	return s.repo.DeleteConsultation(idConsultation)
}

func (s *ConsultationRepo) CreateConsultationNotes(notes Repos.Notes) (Repos.Notes, error) {
	return s.repo.CreateConsultationNotes(notes)
}
func (s *ConsultationRepo) UpdateConsultationNotes(notes Repos.Notes) (Repos.Notes, error) {
	return s.repo.UpdateConsultationNotes(notes)
}

func (s *ConsultationRepo) DeleteConsultationNotes(idNotes int) error {
	return s.repo.DeleteConsultationNotes(idNotes)
}
