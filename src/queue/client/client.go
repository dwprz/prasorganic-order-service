package client

import (
	"time"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/interface/queue"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type QueueImpl struct {
	client *asynq.Client
}

func NewQueue() queue.Client {
	client := asynq.NewClient(asynq.RedisClusterClientOpt{
		Addrs: []string{
			config.Conf.Redis.AddrNode1,
			config.Conf.Redis.AddrNode2,
			config.Conf.Redis.AddrNode3,
			config.Conf.Redis.AddrNode4,
			config.Conf.Redis.AddrNode5,
			config.Conf.Redis.AddrNode6,
		},
		Password: config.Conf.Redis.Password,
	})

	return &QueueImpl{
		client: client,
	}
}

func (q *QueueImpl) Create(typename string, queue string, payload []byte, delay time.Duration) {
	opt := []asynq.Option{
		asynq.ProcessIn(delay),
		asynq.MaxRetry(5),
		asynq.Timeout(time.Duration(20 * time.Second)),
	}

	task := asynq.NewTask(typename, payload, opt...)
	res, err := q.client.Enqueue(task, asynq.Queue(queue))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.Queue/Create", "section": "client.Enqueue"}).Error(err)
		return
	}

	log.Logger.WithFields(logrus.Fields{"location": "client.Queue/Create", "section": "client.Enqueue"}).Infof("successfully enqueued: %+v", res)
}

func (q *QueueImpl) Close() {
	if err := q.client.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.Queue/Create", "section": "client.Enqueue"}).Error(err)
	}
}
