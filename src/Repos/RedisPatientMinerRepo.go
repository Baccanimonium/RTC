package Repos

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	patientMinersList            = "patientMinersList:"
	patientMinersStruct          = "patientMinersStruct:"
	patientMinersTags            = "patientMinersTags:"
	patientCandidatesList        = "patientCandidates:"
	patientCandidatesStruct      = "patientCandidatesStruct:"
	patientCandidatesTags        = "patientCandidatesTags:"
	patientCandidatesLockList    = "patientCandidatesLockList:"
	patientCandidatesReleaseList = "patientCandidatesReleaseList:"
	patientMinersExcludeList     = "patientMinersExcludeList:"
)

type PatientMiner struct {
	IdUser   int     `bson:"id_user" json:"id_user" binding:"required"`
	Tags     []int   `bson:"tags" json:"tags,omitempty"`
	Lat      float64 `bson:"lat" json:"lat"`
	Long     float64 `bson:"long" json:"long"`
	Index    int     `bson:"index" json:"index"`
	QueueKey string  `bson:"queue_key" json:"queue_key"`
}

type RedisPatientCandidate struct {
	IdUser      int     `bson:"id_user" json:"id_user"`
	Tags        []int   `bson:"tags" json:"tags"`
	Lat         float64 `bson:"lat" json:"lat"`
	Long        float64 `bson:"long" json:"long"`
	Description string  `bson:"description" json:"description"`
	CreatedAt   int64   `bson:"created_at" json:"created_at"`
}

type SearchPatientCandidateParams struct {
	ctx                          context.Context
	pipe                         redis.Pipeliner
	patientMiner                 *PatientMiner
	patientCandidatesCache       map[string]RedisPatientCandidate
	minedPatients                map[int]RedisPatientCandidate
	patientCandidatesIds         []string
	startIndex                   int
	endIndex                     int
	minerId                      string
	unixTaskTime                 int64
	releasePatientCandidatesTime string
}

/* TODO возможно стоит хранить распарсенный список пациентов в оперативке(в *r)  и не вычитывать его из редиса, при каждом запросе.
TODO а при каждой таске(раз в минуту, вычитывать его из редиса), тогда индивидуальные запросы будут парсить меньше
TODO писать индивидуальный запрос, пролонгирование, удаления, общий рефактор после этого
*/
func (r *Redis) searchPatientCandidate(searchParams *SearchPatientCandidateParams) (int, bool) {
	// перебираем всех доступных пациентов
	for i := searchParams.startIndex; i < searchParams.endIndex; i++ {
		patientCandidateId := searchParams.patientCandidatesIds[i]
		// если кандидат не залочен и не находиться в списке исключений
		if !searchParams.pipe.SIsMember(searchParams.ctx, patientCandidatesLockList, patientCandidateId).Val() &&
			!searchParams.pipe.SIsMember(searchParams.ctx, patientMinersExcludeList+searchParams.minerId, patientCandidateId).Val() {
			patientCandidate, err := r.GetPatientCandidate(patientCandidateId, searchParams.patientCandidatesCache)
			if err == nil {
				var isTagIntersect bool
				// если кандидат в пациенты встал в очередь более суток назад, мы предлагаем его всем врачам
				if searchParams.unixTaskTime-patientCandidate.CreatedAt > 86400 {
					isTagIntersect = true
					// иначе проверяем его теги и локации
					// TODO: Реализовать логику подсчета координат
				} else {
					for _, tag := range searchParams.patientMiner.Tags {
						if searchParams.pipe.SIsMember(searchParams.ctx, patientCandidatesTags+patientCandidateId, tag).Val() {
							isTagIntersect = true
							break
						}

					}
				}

				if isTagIntersect {
					// лочим кандидата в пациенты
					searchParams.pipe.SAdd(searchParams.ctx, patientCandidatesLockList, patientCandidateId)
					// указываем, когда нужно разлочить пачку кандидатов
					searchParams.pipe.SAdd(searchParams.ctx, patientCandidatesReleaseList+searchParams.releasePatientCandidatesTime, patientCandidateId)
					searchParams.minedPatients[searchParams.patientMiner.IdUser] = patientCandidate
					// останавливаем поиск для доктора
					return i, true
				}
			}
		}
	}
	return 0, false
}

func NewRedisPatientMinerRepo(rdb *redis.Client) *Redis {
	return &Redis{rdb: rdb}
}

func (r *Redis) CreatePatientMiner(patientMiner PatientMiner) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		patientMinerId := strconv.Itoa(patientMiner.IdUser)
		patientMinerStructKey := patientMinersStruct + patientMinerId
		queueUpdateKey := time.Now().Local().Add(time.Minute * time.Duration(5)).Format("15:04")

		pipe.HSet(ctx, patientMinerStructKey, "id_patient", patientMiner.IdUser)
		pipe.HSet(ctx, patientMinerStructKey, "lat", patientMiner.Lat)
		pipe.HSet(ctx, patientMinerStructKey, "long", patientMiner.Long)
		pipe.HSet(ctx, patientMinerStructKey, "index", 0)
		pipe.HSet(ctx, patientMinerStructKey, "queue_key", queueUpdateKey)
		pipe.SAdd(ctx, patientMinersTags+patientMinerId, "tags", patientMiner.Tags)

		pipe.SAdd(ctx, patientMinersList+queueUpdateKey, patientMinerId) // храним список по времени
		return nil
	})

	return err
}

func (r *Redis) DeletePatientMiner(patientMinerId int) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		strPatientMinerId := strconv.Itoa(patientMinerId)
		patientMinerStructKey := patientMinersStruct + strPatientMinerId

		queueKey := pipe.HGet(ctx, patientMinerStructKey, "queue_key").Val()
		pipe.Del(ctx, patientMinerStructKey)                      // удаляем структуру
		pipe.SRem(ctx, queueKey, strPatientMinerId)               // удаляем врача из очереди
		pipe.Del(ctx, patientMinersTags+strPatientMinerId)        // удаляем тэги врача
		pipe.Del(ctx, patientMinersExcludeList+strPatientMinerId) // удаляем отфильтрованных врачом пациентов
		return nil
	})

	return err
}

func (r *Redis) ProlongPatientCandidateOffer(patientMinerId, patientCandidateId int) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		strPatientMinerId := strconv.Itoa(patientMinerId)
		strPatientCandidateId := strconv.Itoa(patientCandidateId)
		patientMinerStructKey := patientMinersStruct + strPatientMinerId

		queueKey := pipe.HGet(ctx, patientMinerStructKey, "queue_key").Val()
		releasePatientCandidatesTime := time.Now().Local().Add(time.Minute * time.Duration(30)).Format("15:04")
		// удаляем из старой очереди на разлок
		pipe.SRem(ctx, patientCandidatesReleaseList+queueKey, strPatientCandidateId)
		// указываем, когда нужно разлочить кандидата
		pipe.SAdd(ctx, patientCandidatesReleaseList+releasePatientCandidatesTime, strPatientCandidateId)

		// удаляем врача из планировщика
		pipe.SAdd(ctx, patientMinersList+queueKey, strPatientMinerId)
		// Планируем следующую раздачу пациентов, для врача
		pipe.SAdd(ctx, patientMinersList+releasePatientCandidatesTime, strPatientMinerId)
		// обновляем ключ следующей таски врача, чтобы знать где его удалять
		pipe.HSet(ctx, patientMinerStructKey, "queue_key", releasePatientCandidatesTime)
		return nil
	})

	return err
}

func (r *Redis) CreatePatientCandidate(patientCandidate RedisPatientCandidate) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		patientCandidateId := strconv.Itoa(patientCandidate.IdUser)

		patientCandidateStructKey := patientCandidatesStruct + patientCandidateId
		pipe.HSet(ctx, patientCandidateStructKey, "id_patient", patientCandidate.IdUser)
		pipe.HSet(ctx, patientCandidateStructKey, "lat", patientCandidate.Lat)
		pipe.HSet(ctx, patientCandidateStructKey, "long", patientCandidate.Long)
		pipe.HSet(ctx, patientCandidateStructKey, "description", patientCandidate.Description)
		pipe.HSet(ctx, patientCandidateStructKey, "created_at", time.Now().Unix())

		pipe.SAdd(ctx, patientCandidatesTags+patientCandidateId, "tags", patientCandidate.Tags)

		return nil
	})
	return err
}

func (r *Redis) DeletePatientCandidate(patientCandidateId int) error {
	ctx := context.Background()
	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		strPatientCandidateId := strconv.Itoa(patientCandidateId)
		pipe.Del(ctx, patientCandidatesStruct+strPatientCandidateId)
		pipe.Del(ctx, patientCandidatesTags+strPatientCandidateId)

		return nil
	})
	return err
}

func (r *Redis) GetPatientCandidate(patientCandidateId string, cache map[string]RedisPatientCandidate) (RedisPatientCandidate, error) {
	if val, ok := cache[patientCandidateId]; ok {
		return val, nil
	}

	ctx := context.Background()

	var patientCandidate RedisPatientCandidate
	err := r.rdb.HGetAll(ctx, patientCandidatesStruct+patientCandidateId).Scan(&patientCandidate)
	if err != nil {
		return patientCandidate, err
	}

	err = r.rdb.SMembers(ctx, patientCandidatesTags+patientCandidateId).ScanSlice(&patientCandidate.Tags)

	return patientCandidate, err
}

func (r *Redis) MinePatients(atTime time.Time) (map[int]RedisPatientCandidate, error) {
	taskTime := atTime.Format("15:04")
	releasePatientCandidatesTime := atTime.Add(time.Minute * time.Duration(5)).Format("15:04")
	unixTaskTime := atTime.Unix()
	ctx := context.Background()
	var patientCandidatesRelease []string                       // список кандидатов в пациенты, который нужно разлочить
	var patientCandidatesIds []string                           // список ID кандидатов в пациенты по которому бежим в цикле
	var patientMinersIds []string                               // список ID Докторов которым нужно выдать пациентов
	var patientCandidatesCache map[string]RedisPatientCandidate // кэш прочитанных пациентов из редиса
	var minedPatients map[int]RedisPatientCandidate             // key patientMinerID

	_, err := r.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// читаем из редиса список для анлока
		err := pipe.SMembers(ctx, patientCandidatesReleaseList+taskTime).ScanSlice(&patientCandidatesRelease)
		// удаляем анлок лист
		pipe.Del(ctx, patientCandidatesReleaseList+taskTime)

		if err != nil {
			return err
		}
		// снимаем лок с пациентов, для последующих циклов
		pipe.SRem(ctx, patientCandidatesLockList, patientCandidatesRelease)

		// получаем список пациентов
		err = pipe.SMembers(ctx, patientCandidatesList).ScanSlice(&patientCandidatesIds)
		if err != nil {
			return err
		}
		patientCandidatesIdsLength := len(patientCandidatesIds)

		// получаем список докторов, которым нужно раздать новых пациентов
		err = pipe.SMembers(ctx, patientMinersList+taskTime).ScanSlice(&patientMinersIds)
		// удаляем список докторов
		pipe.Del(ctx, patientMinersList+taskTime)
		if err != nil {
			return err
		}

		// ищем каждому врачу нового пациента
		for _, minerId := range patientMinersIds {
			patientMinerStructKey := patientMinersStruct + minerId

			var patientMiner PatientMiner
			err = pipe.HGetAll(ctx, patientMinerStructKey).Scan(&patientMiner)
			if err != nil {
				return err
			}

			err := pipe.SMembers(ctx, patientMinersTags+minerId).ScanSlice(&patientMiner.Tags)
			if err != nil {
				return err
			}

			searchParams := SearchPatientCandidateParams{
				ctx:                          ctx,
				pipe:                         pipe,
				patientMiner:                 &patientMiner,
				patientCandidatesCache:       patientCandidatesCache,
				minedPatients:                minedPatients,
				patientCandidatesIds:         patientCandidatesIds,
				startIndex:                   patientMiner.Index,
				endIndex:                     patientCandidatesIdsLength,
				minerId:                      minerId,
				unixTaskTime:                 unixTaskTime,
				releasePatientCandidatesTime: releasePatientCandidatesTime,
			}

			// ищем пациента для врача
			nextIndex, success := r.searchPatientCandidate(&searchParams)

			// если не находим, начинаем поиск сначала
			if !success {
				searchParams.startIndex = 0
				searchParams.endIndex = patientMiner.Index
				index, _ := r.searchPatientCandidate(&searchParams)
				nextIndex = index
			}

			// перебираем всех доступных пациентов
			//for i := patientMiner.Index; i < patientCandidatesIdsLength; i++ {
			//	patientCandidateId := patientCandidatesIds[i]
			//	// если кандидат не залочен и не находиться в списке исключений
			//	if !pipe.SIsMember(ctx, patientCandidatesLockList, patientCandidateId).Val() &&
			//		!pipe.SIsMember(ctx, patientMinersExcludeList+minerId, patientCandidateId).Val() {
			//		patientCandidate, err := r.GetPatientCandidate(patientCandidateId, patientCandidatesCache)
			//		if err == nil {
			//			var isTagIntersect bool
			//			// если кандидат в пациенты встал в очередь более суток назад, мы предлагаем его всем врачам
			//			if unixTaskTime-patientCandidate.CreatedAt > 86400 {
			//				isTagIntersect = true
			//				// иначе проверяем его теги и локации
			//				// TODO: Реализовать логику подсчета координат
			//			} else {
			//				for _, tag := range patientMiner.Tags {
			//					if pipe.SIsMember(ctx, patientCandidatesTags+patientCandidateId, tag).Val() {
			//						isTagIntersect = true
			//						break
			//					}
			//
			//				}
			//			}
			//
			//			if isTagIntersect {
			//				// лочим кандидата в пациенты
			//				pipe.SAdd(ctx, patientCandidatesLockList, patientCandidateId)
			//				// указываем, когда нужно разлочить пачку кандидатов
			//				pipe.SAdd(ctx, patientCandidatesReleaseList+releasePatientCandidatesTime, patientCandidateId)
			//				minedPatients[patientMiner.IdUser] = patientCandidate
			//				// останавливаем поиск для доктора
			//				break
			//			}
			//		}
			//	}
			//}

			// Планируем следующую раздачу пациентов, для врача
			pipe.SAdd(ctx, patientMinersList+releasePatientCandidatesTime, minerId)
			// обновляем ключ следующей таски врача, чтобы знать где его удалять
			pipe.HSet(ctx, patientMinerStructKey, "queue_key", releasePatientCandidatesTime)
			// обновляем индекс текущего пациента, чтобы знать, откуда начинать следующий поиск
			pipe.HSet(ctx, patientMinerStructKey, "index", nextIndex)

		}

		return nil
	})

	return minedPatients, err
}
