package Services

import (
	"video-chat-app/src/Models"
	"video-chat-app/src/Repos"
)

type PatientCandidates struct {
	repo Repos.PatientCandidatesRepo
}

func NewPatientCandidatesService(repo Repos.PatientCandidatesRepo) *PatientCandidates {
	return &PatientCandidates{repo: repo}
}

func (s *PatientCandidates) CreatePatientCandidate(patientCandidate Models.PatientCandidate) (interface{}, error) {
	return s.repo.CreatePatientCandidate(patientCandidate)
}

func (s *PatientCandidates) GetAllPatientCandidates() ([]Models.PatientCandidate, error) {
	return s.repo.GetAllPatientCandidates()
}
