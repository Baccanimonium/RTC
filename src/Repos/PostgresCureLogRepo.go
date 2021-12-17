package Repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

/*
	Прикрепленный текст и файлы должны по выполнению задания должны попасть сюда с поменткой
	выполнение задания "текст задания" и типом taskComplete
	Момент создания задания должен быть добавлен как новое задание "текст задания", "описание задания" и тип newTask
	По завершению консультации "Заметки консультации" должны попадать сюда с заголовком "Заметки консультации" и
	типом consultationsNotes

	поиск по типу
	поиск по датам
*/

type LogItem struct {
	Id        int    `json:"id,omitempty"`
	IdPatient int    `json:"id_patient,omitempty" binding:"required"`
	Text      string `json:"text,omitempty"`
	File      string `json:"file,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	LogType   string `json:"log_type,omitempty"`
}

func NewCureLogPostgresRepo(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateLogItem(logItem LogItem) (LogItem, error) {
	var newLogItem LogItem

	query := fmt.Sprintf(
		"INSERT INTO %s (id_patient, text, file, log_type, created_at) values ($1, $2, $3, $4, $5) RETURNING *",
		patientsLog,
	)
	row := r.db.QueryRow(query, logItem.IdPatient, logItem.Text, logItem.File, logItem.LogType, time.Now().Unix())

	if err := row.Scan(&logItem); err != nil {
		return LogItem{}, err
	}

	return newLogItem, nil
}

type LogItemSearchParams struct {
	LogType  string `json:"log_type"`
	FromDate int64  `json:"from_date"`
	TillDate int64  `json:"till_date"`
}

func (r *Postgres) GetAllLogItem(params LogItemSearchParams, idPatient int) ([]LogItem, error) {
	var logItems []LogItem
	var conditionQuery []string
	argsArray := make([]interface{}, 0)
	argsCount := 0

	argsCount += 1
	conditionQuery = append(conditionQuery, fmt.Sprintf(
		"id_patient = $%d",
		argsCount,
	))

	argsArray = append(argsArray, idPatient)

	if params.LogType != "" {
		argsCount += 1
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_patient = $%d",
			argsCount,
		))
		argsArray = append(argsArray, params.LogType)
	}

	if params.FromDate != 0 {
		argsCount += 1
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"from_date > $%d",
			argsCount,
		))

		argsArray = append(argsArray, params.FromDate)
	}

	if params.TillDate != 0 {
		argsCount += 1
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"till_date < $%d",
			argsCount,
		))

		argsArray = append(argsArray, params.TillDate)
	}

	query := fmt.Sprintf(`SELECT * FROM %s doc WHERE %s`, patientsLog, strings.Join(conditionQuery, " AND "))

	err := r.db.Select(&logItems, query, argsArray...)

	return logItems, err
}

func (r *Postgres) GetLogItemById(id int) (LogItem, error) {
	var logItem LogItem

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, patientsLog)

	err := r.db.Get(&logItem, query, id)

	return logItem, err
}

func (r *Postgres) UpdateLogItemBy(logItem LogItem) (LogItem, error) {
	var newLogItem LogItem

	query := fmt.Sprintf(
		"UPDATE %s SET id_patient=$1, text=$2, file=$3 WHERE id=$4 RETURNING *",
		patientsLog,
	)
	err := r.db.Get(&newLogItem, query, logItem.IdPatient, logItem.Text, logItem.File, logItem.Id)

	return newLogItem, err
}

func (r *Postgres) DeleteLogItem(id int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id=$1",
		patientsLog,
	)
	_, err := r.db.Exec(query, id)

	return err
}
