package Tasks

import "github.com/sirupsen/logrus"

/*
	1. Пациенты вставшие в очередь хранятся в ней до тех пор, пока доктор не примет их как пациента
	2. Врачи вставшие в очередь хранятся в ней пока они онлайн
	3. Если какой-то врач смотрит на пациента. То на него не может смотреть ни какой врач
	4. Если врач смотрел на пациента, то он не должен выдаваться ему повторно
	5. Поиск должен быть оптимизирован по тэгам и геолокации. Предпочитаемый радиус выставляет пользователь
	6. Поиск должен помнить индекс пациента, которого смотрит доктор и начинать повторный поиск с него
	7. Поиск должен помнить номер круга поиска врача, чтобы выбирать пациентов без учета геолокации и тегов, при условии
		что, пациент в очереди более 1 дня.
*/

func (tm *TaskManager) RunPatientMinerTask() error {
	_, err := tm.scheduler.Every(1).Minute().Tag("PatientMiner").Do(tm.eventTask)
	return err
}

func (tm *TaskManager) PatientMinerTask() {
	patientCandidates, err := tm.services.PatientCandidatesService.GetAllPatientCandidates()
	if err != nil {
		logrus.Print("Failed To get patient Candidates list, ERROR ", err.Error())
	}

	for _, event := range events {

	}
}
