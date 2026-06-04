package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FiyZou/handygo/queue"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const TypeUserReport = "user:report"

type UserReportPayload struct {
	UserID uint `json:"userId"`
}

func NewUserReportTask(userID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(UserReportPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeUserReport, payload), nil
}

func Register(server *queue.Server, logger *zap.Logger) {
	server.HandleFunc(TypeUserReport, handleUserReport(logger))
}

func handleUserReport(logger *zap.Logger) func(context.Context, *asynq.Task) error {
	if logger == nil {
		logger = zap.NewNop()
	}
	return func(ctx context.Context, task *asynq.Task) error {
		var payload UserReportPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("decode user report payload: %w", err)
		}
		logger.Info("handled user report task", zap.Uint("user_id", payload.UserID))
		return nil
	}
}
