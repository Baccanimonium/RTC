package Models

type DoctorSchedule struct {
	IdDoctor    int      `json:"id_doctor"`
	WorkPeriods []Period `json:"work_periods"`
}

type Period struct {
	Start int64 `json:"condition"`
	End   int64 `json:"end"`
	Day   int   `json:"day"`
}

type GetDoctorSchedule struct {
	IdDoctor int   `json:"id_doctor"`
	Date     int64 `json:"date"`
}
