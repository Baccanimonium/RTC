package Models

type Consultation struct {
	Id           int    `json:"id"`
	IdPatient    int    `json:"id_patient" binding:"required"`
	IdDoctor     int    `json:"id_doctor"  binding:"required"`
	Start        int64  `json:"start" binding:"required"`
	End          int64  `json:"end"       binding:"required"`
	Offline      bool   `json:"offline"    binding:"required"`
	DoctorJoined string `json:"doctor_joined"`
}

type Notes struct {
	Id             int      `json:"id"`
	Notes          int      `json:"notes" binding:"required"`
	IdConsultation int      `json:"id_consultation" binding:"required"`
	ForDoctor      bool     `json:"for_doctor"`
	Files          []string `json:"files"`
}

type GetConsultationList struct {
	IdDoctor  int
	IdPatient int
	Start     int64
	End       int64
}
