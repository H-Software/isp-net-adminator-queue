package queue

import (
	"context"
	"github.com/h-software/isp-net-adminator-queue/internal/flag"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
	"github.com/hibiken/asynq"
	"strings"
	"time"
)

const (
	TypeAdminatorWorkItem = "adminator3:workitem:3" // adminator3:workitem:basic

	TypeAdminatorWorkItemAgg = "adminator3:workitem:agg"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.NewLogger(nil)
}

func aggregate(group string, tasks []*asynq.Task) *asynq.Task {
	logger.Infof("Aggregating %d tasks from group %q", len(tasks), group)
	var b strings.Builder
	for _, t := range tasks {
		b.Write(t.Payload())
		b.WriteString("\n")
	}
	return asynq.NewTask(TypeAdminatorWorkItemAgg, []byte(b.String()))
}

func HandleWorkItemAggTask(ctx context.Context, task *asynq.Task) error {
	logger.Infof("Handler received aggregated task")
	logger.Infof("aggregated messages: %s", task.Payload())

	//     if err := json.Unmarshal(t.Payload(), &p); err != nil {
	//     return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	// }
	// logger.Infof("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)

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
