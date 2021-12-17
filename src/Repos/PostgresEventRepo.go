package Repos

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Id                   int    `json:"id"`
	IdPatient            int    `json:"id_patient" binding:"required"`
	IdDoctor             int    `json:"id_doctor" binding:"required"`
	AtDays               []Day  `json:"at_days"`
	Title                string `json:"title"`
	Description          string `json:"description"`
	NotifyDoctor         bool   `json:"notify_doctor"`
	RemindInAdvance      int    `json:"remind_in_advance"` // В секундах
	ConfirmationTime     int    `json:"confirmation_time"` // В секундах
	RequiresConfirmation bool   `json:"requires_confirmation"`
	CreatedAt            string `json:"created_at"`
	LastTill             string `json:"last_till"`
	Weight               int    `json:"weight"`
}

type Day struct {
	Condition int    `json:"condition"`
	Day       int    `json:"day"`
	Result    int    `json:"result"`
	AtDay     string `json:"at_day"`
	AtTime    string `json:"at_time"`
}

func (c Day) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Day) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {

		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func (r *EventPostgres) getCurrentWeek(date time.Time) (int, error) {
	// текущую неделю считаем от нулевой недели в 1970
	// 604800 - неделя в секундах
	currentWeek := int(date.In(r.loc).Unix() / 604800)

	return currentWeek, nil
}

func (r *EventPostgres) prepareDayToSave(day Day, currentWeek, RemindInAdvance int) (Day, error) {

	// определяем по каким неделям у нас происходит событие, currentWeek - текуща	я неделя от нулевой недели в
	// 1970 году. Condition - условие к примеру 2147 % 4 = 3
	day.Result = currentWeek % day.Condition
	eventTime, err := time.Parse("2006-01-02T15:04Z07:00", "1970-01-01T"+day.AtTime)
	if err != nil {
		return day, err
	}
	day.Day = r.dayDictionary[day.AtDay] + int(eventTime.In(r.loc).Unix()) - RemindInAdvance
	//
	isDateWeekDay := false

	for _, v := range r.weekDayDictionary {
		if v == day.AtDay {
			isDateWeekDay = true
		}
	}

	var valueOfTheLastDay int
	var valueOfTheFirstDay int

	if isDateWeekDay {
		valueOfTheLastDay = r.dayDictionary["Su"]
		valueOfTheFirstDay = r.dayDictionary["Mo"]
	} else {
		valueOfTheLastDay = r.dayDictionary["31"]
		valueOfTheFirstDay = r.dayDictionary["1"]
	}

	if day.Day < valueOfTheFirstDay {
		day.Day = valueOfTheLastDay + 86400 + day.Day - valueOfTheFirstDay
	}

	return day, nil
}

func (r *EventPostgres) prepareEventToSave(event Event) ([]Day, error) {
	// время запаковываем в day каждого элемента
	// remind in advance просчитываем при создании условий в at_days, храним уведомления
	// как дополнительные условия, туда же просчитываем результат деления недели
	var result = make([]Day, 0)
	// Получаем дату, находим тайм зону, приводим ее к нулевой там зоне
	//"2021-11-01T13:29:05+07:00" пример даты
	t, err := time.Parse(time.RFC3339, event.CreatedAt)

	if err != nil {
		return result, err
	}

	currentWeek, err := r.getCurrentWeek(t)

	if err != nil {
		return result, err
	}

	for _, day := range event.AtDays {
		updatedDay, err := r.prepareDayToSave(day, currentWeek, 0)
		if err != nil {
			return result, err
		}
		result = append(result, updatedDay)

		if event.RemindInAdvance != 0 {
			dayInAdvance, err := r.prepareDayToSave(day, currentWeek, event.RemindInAdvance)
			if err != nil {
				return result, err
			}
			result = append(result, dayInAdvance)
		}
	}

	return result, nil
}

/* TODO добавить время подтверждения 3600 секунд на пример,
   TODO пред расчет дат уведомления заранее, то есть уведомляем за
   TODO 3600 секунд на подтверждения эв 90000 секунд(1 день и 1 час)
   TODO значит atTimeConfirmation на 1 час меньше чем AtTime
   TODO AtDate на 1 день меньше, Days тоже смещаются на 1
*/

type EventPostgres struct {
	db                *sqlx.DB
	loc               *time.Location
	dayDictionary     map[string]int
	weekDayDictionary [7]string
	dateDictionary    [31]string
}

func NewEventPostgresRepo(db *sqlx.DB) *EventPostgres {
	loc, err := time.LoadLocation("Africa/Accra")
	if err != nil {
		logrus.Fatal(err.Error())
	}

	dayDictionary := map[string]int{
		"Mo": 345600,
		"Tu": 432000,
		"We": 518400,
		"Th": 604800,
		"Fr": 691200,
		"Sa": 777600,
		"Su": 864000,
		"1":  2678400,
		"2":  2764800,
		"3":  2851200,
		"4":  2937600,
		"5":  3024000,
		"6":  3110400,
		"7":  3196800,
		"8":  3283200,
		"9":  3369600,
		"10": 3456000,
		"11": 3542400,
		"12": 3628800,
		"13": 3715200,
		"14": 3801600,
		"15": 3888000,
		"16": 3974400,
		"17": 4060800,
		"18": 4147200,
		"19": 4233600,
		"20": 4320000,
		"21": 4406400,
		"22": 4492800,
		"23": 4579200,
		"24": 4665600,
		"25": 4752000,
		"26": 4838400,
		"27": 4924800,
		"28": 5011200,
		"29": 5097600,
		"30": 5184000,
		"31": 5270400,
	}
	weekDayDictionary := [7]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	dateDictionary := [31]string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
		"20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31",
	}

	return &EventPostgres{
		db:                db,
		loc:               loc,
		dayDictionary:     dayDictionary,
		weekDayDictionary: weekDayDictionary,
		dateDictionary:    dateDictionary,
	}
}

func (r *EventPostgres) CreateEvent(event Event) (Event, error) {
	var result Event

	query := fmt.Sprintf(
		`INSERT INTO %s
                (id_patient, id_doctor, at_days, title, description, notify_doctor, remind_in_advance,
				confirmation_time, requires_confirmation, created_at, last_till, weight)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				RETURNING *`,
		eventTable,
	)

	days, err := r.prepareEventToSave(event)

	if err != nil {
		return Event{}, err
	}

	event.AtDays = days

	row := r.db.QueryRow(
		query,
		event.IdPatient,
		event.IdDoctor,
		event.AtDays,
		event.Title,
		event.Description,
		event.NotifyDoctor,
		event.RemindInAdvance,
		event.ConfirmationTime,
		event.RequiresConfirmation,
		event.CreatedAt,
		event.LastTill,
		event.Weight,
	)

	err = row.Scan(&result)

	return result, err
}

func (r *EventPostgres) UpdateEvent(event Event) (Event, error) {
	var result Event

	query := fmt.Sprintf(
		`UPDATE %s
				SET at_days=$1, title=$2, description=$3, notify_doctor=$4, remind_in_advance=$5,
				confirmation_time=$6, requires_confirmation=$7, created_at=$8, last_till=$9, weight=$10
				WHERE id=$11 RETURNING *`,
		eventTable,
	)

	days, err := r.prepareEventToSave(event)

	if err != nil {
		return Event{}, err
	}

	event.AtDays = days

	err = r.db.Get(
		&result,
		query,
		event.AtDays,
		event.Title,
		event.Description,
		event.NotifyDoctor,
		event.RemindInAdvance,
		event.ConfirmationTime,
		event.RequiresConfirmation,
		event.CreatedAt,
		event.LastTill,
		event.Weight,
		event.Id,
	)

	return result, err
}

func (r *EventPostgres) GetEventById(id int) (Event, error) {
	var event Event

	query := fmt.Sprintf("SELECT * %s WHERE id = $1", eventTable)

	err := r.db.Get(&event, query, id)

	return event, err
}

type GetAllEventsParams struct {
	IdPatient int
	IdDoctor  int
	Date      time.Time
}

func (r *EventPostgres) GetAllEvents(request GetAllEventsParams) ([]Event, error) {
	var result = make([]Event, 0)
	var conditionQuery []string
	argsCount := 0
	argsArray := make([]interface{}, 0)

	if request.IdPatient != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_patient = $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, request.IdPatient)
	}

	if request.IdDoctor != 0 {
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			"id_doctor = $%d",
			argsCount+1,
		))
		argsCount += 1
		argsArray = append(argsArray, request.IdDoctor)
	}

	if !request.Date.IsZero() {
		currentWeek, err := r.getCurrentWeek(request.Date)
		if err != nil {
			return result, err
		}
		dateTime := int(time.Date(1970, 1, 1, request.Date.Hour(), request.Date.Minute(), 0, 0, r.loc).Unix())
		currentWeekday := request.Date.Weekday().String()[0:2]
		currentDay := strconv.Itoa(request.Date.Day())
		conditionQuery = append(conditionQuery, fmt.Sprintf(
			`arr.item_object->>'day'= ANY($%d)
			and $%d %% cast(arr.item_object->>'condition' as int) = cast(arr.item_object->>'result' as int)`,
			argsCount+1,
			argsCount+2,
		))
		argsCount += 2
		argsArray = append(
			argsArray,
			[3]int{
				r.dayDictionary[currentWeekday] + dateTime,
				r.dayDictionary[currentDay] + dateTime,
				int(request.Date.Unix()),
			},
			currentWeek,
		)
	}

	query := fmt.Sprintf(
		"SELECT * FROM, jsonb_array_elements(at_days) with ordinality arr(item_object, position) %s"+" WHERE %s",
		eventTable,
		strings.Join(conditionQuery, " AND "),
	)

	err := r.db.Select(&result, query, argsArray...)

	return result, err
}

func (r *EventPostgres) DeleteEvent(id int) (Event, error) {
	var result Event

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING *", eventTable)
	err := r.db.Get(&result, query, id)

	return result, err
}
