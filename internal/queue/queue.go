package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/h-software/isp-net-adminator-queue/internal/flag"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
	"github.com/hibiken/asynq"
)

const (
	TypeAdminatorWorkItem = "adminator3:workitem:3" // adminator3:workitem:basic

	TypeAdminatorWorkItemAgg = "adminator3:workitem:agg"

	payloadSeparator = "\n"
)

type AdminatorTaskType int

const (
	WorkItem AdminatorTaskType = iota
	EmailItem
)

var (
	logger *log.Logger
)

type WorkItemPayload struct {
	ItemId    int    `json:"item_id"`
	CreatedAt int    `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
}

func init() {
	logger = log.NewLogger(nil)
}

func aggregate(group string, tasks []*asynq.Task) *asynq.Task {
	logger.Infof("Aggregating %d tasks from group %q", len(tasks), group)
	var b strings.Builder
	for _, t := range tasks {
		b.Write(t.Payload())
		b.WriteString(payloadSeparator)
	}

	return asynq.NewTask(TypeAdminatorWorkItemAgg, []byte(b.String()))
}

func HandleWorkItemAggTask(ctx context.Context, task *asynq.Task) error {
	logger.Infof("Handler received aggregated task")
	// logger.Infof("aggregated messages: %s", task.Payload())

	id, ok := asynq.GetTaskID(ctx)
	if !ok {
		logger.Errorf("get TaskID for aggregated message failed")
	} else {
		logger.Infof("aggregated message TaskID: %s", id)
	}

	if err := HandleAggTaskPayload(ctx, task, id, WorkItem); err != nil {
		errM := fmt.Sprintf("handle aggreaged task payload failed: %v, id: %s", err, id)
		logger.Error(errM)
		return fmt.Errorf("%s", errM)
	}

	return nil
}

func HandleAggTaskPayload(ctx context.Context, task *asynq.Task, taskId string, taskType AdminatorTaskType) error {
	var payloadParsed WorkItemPayload
	itemIdChecksum := 0
	itemIdCount := 0

	// parse payload
	payloadMap := bytes.Split(task.Payload(), []byte(payloadSeparator))

	for _, payload := range payloadMap {
		// Split also saves empty lines
		if len(string(payload)) > 0 {
			logger.Infof("parsed payload: %s", string(payload))

			if err := json.Unmarshal(payload, &payloadParsed); err != nil {
				errM := fmt.Sprintf("json.Unmarshal failed: %v, id: %s, payload: %s", err, taskId, string(payload))
				logger.Error(errM)
				return fmt.Errorf("%s", errM)
			} else {
				logger.Infof("json.Unmarshal OK")
				logger.Infof("Unmarshaled payload: id=%d createdAt=%d, createdBy=%s", payloadParsed.ItemId, payloadParsed.CreatedAt, payloadParsed.CreatedBy)

				itemIdChecksum += payloadParsed.ItemId
				itemIdCount += 1
			}
		}
	}

	// checksum for itemId
	if taskType == WorkItem {

		check := math.Mod(float64(itemIdChecksum), float64(itemIdCount))

		if check == 0 {
			logger.Infof("checksum for ItemId OK\n")
		} else {
			logger.Infof("checksum for ItemId failed\n")
		}
	} else {
		logger.Infof("checksum for ItemId failed\n")
	}

	return nil
}

func RunServer() *asynq.Server {

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:        *flag.FlagRedisAddr,
			DialTimeout: 2 * time.Second,
		},
		asynq.Config{
			Logger: logger,
			Queues: map[string]int{
				"adminator3:workitem": 3,
			},
			Concurrency:      1,
			GroupAggregator:  asynq.GroupAggregatorFunc(aggregate),
			GroupGracePeriod: *flag.FlagGroupGracePeriod,
			GroupMaxDelay:    *flag.FlagGroupMaxDelay,
			GroupMaxSize:     *flag.FlagGroupMaxSize,
		},
	)

	return srv
}
