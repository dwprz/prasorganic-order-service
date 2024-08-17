package queue

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)

	log.Logger.WithFields(logrus.Fields{
		"location": "asynq.ErrorHandler",
		"taskType": task.Type,
		"retried":  retried,
		"maxRetry": maxRetry,
	}).Error(err)
}
