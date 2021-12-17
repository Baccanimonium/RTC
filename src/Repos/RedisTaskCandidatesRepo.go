package Repos

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	redisTaskCandidatesStruct        = "taskCandidatesStruct:"
	redisTaskCandidatesList          = "taskCandidatesList:"
	redisTaskCandidatesListByPatient = "taskCandidatesListByPatient:"
)

type TaskCandidate struct {
	IdPatient int    `json:"id_patient" binding:"required"`
	IdDoctor  int    `json:"id_doctor"  binding:"required"`
	IdEvent   int    `json:"id_event"   binding:"required"`
	Title     string `json:"title"`
	Time      string `json:"at_time"`
	Date      string `json:"at_date"`
	Weight    int    `json:"weight"`
}

func NewTaskCandidatesRedisRepo(rdb *redis.Client) *Redis {
	return &Redis{rdb: rdb}
}

func (r *Redis) CreateTaskCandidates(taskCandidates []Event, currentDate, currentTime string) error {
	ctx := context.Background()
	taskCandidatesKey := redisTaskCandidatesList + currentTime
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, candidate := range taskCandidates {
			candidateId := strconv.Itoa(candidate.Id)
			eventIdWithPrefix := redisTaskCandidatesStruct + candidateId

			pipe.HSet(ctx, eventIdWithPrefix, "id_event", candidate.Id)
			pipe.HSet(ctx, eventIdWithPrefix, "id_patient", candidate.IdPatient)
			pipe.HSet(ctx, eventIdWithPrefix, "id_doctor", candidate.IdDoctor)
			pipe.HSet(ctx, eventIdWithPrefix, "title", candidate.Title)
			// приводим к стандартному формату даты "2021-11-01T13:29:05+07:00"
			pipe.HSet(ctx, eventIdWithPrefix, "at_date", currentDate+"T"+currentTime+"+00:00")
			pipe.HSet(ctx, eventIdWithPrefix, "at_time", currentTime)
			pipe.HSet(ctx, eventIdWithPrefix, "weight", candidate.Weight)

			pipe.SAdd(ctx, taskCandidatesKey, candidateId) // храним список по времени

			patientId := strconv.Itoa(candidate.IdPatient)
			patientIdWithPrefix := redisTaskCandidatesListByPatient + patientId
			pipe.SAdd(ctx, patientIdWithPrefix, candidateId) // храним список по юзеру
		}
		return nil
	})
	return err
}

func (r *Redis) DeleteTaskCandidate(candidateId int) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		eventId := strconv.Itoa(candidateId)
		candidateKey := redisTaskCandidatesStruct + eventId

		var candidate TaskCandidate
		err := r.rdb.HGetAll(ctx, candidateKey).Scan(&candidate)

		if err != nil {
			return err
		}

		unConfirmedEventListKey := redisTaskCandidatesList + candidate.Time

		patientId := strconv.Itoa(candidate.IdPatient)
		patientIdWithPrefix := redisTaskCandidatesListByPatient + patientId

		// удаляем эвент у кандидата из списка не подтвержденных эвентов
		pipe.SRem(ctx, patientIdWithPrefix, candidateKey)
		pipe.Del(ctx, candidateKey)
		pipe.SRem(ctx, unConfirmedEventListKey, eventId)

		return nil
	})
	return err
}

func (r *Redis) ExtractTaskCandidates(taskTime string) ([]TaskCandidate, error) {
	var taskCandidatesIds []string
	var taskCandidates []TaskCandidate
	ctx := context.Background()
	taskCandidatesKey := redisTaskCandidatesList + taskTime

	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := r.rdb.SMembers(ctx, taskCandidatesKey).ScanSlice(&taskCandidatesIds)
		r.rdb.Del(ctx, taskCandidatesKey) // забываем о списке кандидатов после его получения

		if len(taskCandidatesIds) > 0 && err == nil {
			for _, candidateId := range taskCandidatesIds {
				var candidate TaskCandidate
				candidateKey := redisTaskCandidatesStruct + candidateId
				err = r.rdb.HGetAll(ctx, candidateKey).Scan(&candidate)
				r.rdb.Del(ctx, candidateKey) // забываем о кандидате после его получения

				patientId := strconv.Itoa(candidate.IdPatient)
				patientIdWithPrefix := redisTaskCandidatesListByPatient + patientId

				// удаляем эвент у кандидата из списка не подтвержденных эвентов
				pipe.SRem(ctx, patientIdWithPrefix, candidateId)

				if err != nil {
					logrus.Print("FAIL TO READ TASK CANDIDATE ", err.Error())
					break
				}
				taskCandidates = append(taskCandidates, candidate)
			}
		}

		return err
	})

	return taskCandidates, err
}

func (r *Redis) GetTaskCandidatesByPatient(patientId int) ([]TaskCandidate, error) {
	var taskCandidatesIds []string
	var taskCandidates []TaskCandidate
	ctx := context.Background()
	patientIdString := strconv.Itoa(patientId)
	patientTaskCandidatesKey := redisTaskCandidatesListByPatient + patientIdString

	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := r.rdb.SMembers(ctx, patientTaskCandidatesKey).ScanSlice(&taskCandidatesIds)

		if len(taskCandidatesIds) > 0 && err == nil {
			for _, candidateId := range taskCandidatesIds {
				var candidate TaskCandidate
				candidateKey := redisTaskCandidatesStruct + candidateId
				err = r.rdb.HGetAll(ctx, candidateKey).Scan(&candidate)

				if err != nil {
					logrus.Print("FAIL TO READ TASK CANDIDATE BY PATIENT ", err.Error())
					break
				}
				taskCandidates = append(taskCandidates, candidate)
			}
		}

		return err
	})

	return taskCandidates, err
}
