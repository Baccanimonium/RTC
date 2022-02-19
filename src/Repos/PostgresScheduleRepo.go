package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"video-chat-app/src/Models"
)

/*
	1. вариант
		1.1 Пациент ищет врачей по тэгам
		1.2 Находит нужного ищет запись

	2. вариант
		1.1 Пациент ищет врачей по тэгам и дате на которую он хочет попасть
		SELECT всех врачей у которых в тэгах есть метки JOIN расписание, где день соовествуюет границам

	Вывод, нужно апи которое выдаст консультаци и расписание по дню

	Решение, для поддержания гранулярности данных
		1. расписания джойнятся на профиль доктора
		2. консультации селектяться по дню отдельно
		3. фронт вычитает из рабочих периодов созданные консультации и получает доступное время
*/

func NewSchedulePostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error) {
	var result Models.DoctorSchedule

	query := fmt.Sprintf(
		`INSERT INTO %s (id_doctor, work_periods) values ($1, $2) RETURNING *`,
		scheduleTable,
	)
	row := r.db.QueryRow(
		query,
		schedule.IdDoctor,
		schedule.WorkPeriods,
	)

	if err := row.Scan(&result); err != nil {
		return Models.DoctorSchedule{}, err
	}

	return result, nil
}

func (r *Postgres) UpdateSchedule(schedule Models.DoctorSchedule) (Models.DoctorSchedule, error) {
	var newSchedule Models.DoctorSchedule

	query := fmt.Sprintf(
		"UPDATE %s SET work_periods=$2 WHERE id=$1 RETURNING *",
		scheduleTable,
	)

	err := r.db.Get(&newSchedule, query, schedule.IdDoctor, schedule.WorkPeriods)

	return newSchedule, err
}

func (r *Postgres) GetScheduleByDoctorId(id int) (Models.DoctorSchedule, error) {
	var schedule Models.DoctorSchedule

	query := fmt.Sprintf(
		`SELECT id_doctor, work_periods FROM %s WHERE id = $1`,
		scheduleTable,
	)

	err := r.db.Get(&schedule, query, id)

	return schedule, err
}

func (r *Postgres) GetAllSchedule(params Models.PostgresPagination) ([]Models.DoctorSchedule, error) {
	var schedules = make([]Models.DoctorSchedule, 0)

	query := fmt.Sprintf(
		`SELECT id_doctor, work_periods FROM %s LIMIT $1 OFFSET $2`,
		scheduleTable,
	)

	err := r.db.Select(&schedules, query, params.Limit, params.Skip)

	return schedules, err
}

func (r *Postgres) DeleteSchedule(idDoctor int) (Models.DoctorSchedule, error) {
	var deletedSchedule Models.DoctorSchedule
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id_doctor=$1 RETURNING *",
		scheduleTable,
	)

	err := r.db.Get(&deletedSchedule, query, idDoctor)

	return deletedSchedule, err
}
