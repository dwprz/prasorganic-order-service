package server

import (
	"github.com/dwprz/prasorganic-order-service/src/common/errors/queue"
	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/queue/handler"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewQueue(oh *handler.OrderQueue) *Queue {
	redisConnOpt := asynq.RedisClusterClientOpt{
		Addrs: []string{
			config.Conf.Redis.AddrNode1,
			config.Conf.Redis.AddrNode2,
			config.Conf.Redis.AddrNode3,
			config.Conf.Redis.AddrNode4,
			config.Conf.Redis.AddrNode5,
			config.Conf.Redis.AddrNode6,
		},
		Password: config.Conf.Redis.Password,
	}

	srv := asynq.NewServer(redisConnOpt, asynq.Config{
		Concurrency: 20,
		Queues: map[string]int{
			"orders": 1,
		},
		Logger:   log.Logger,
		ErrorHandler: asynq.ErrorHandlerFunc(queue.ErrorHandler),
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("orders:shipping", oh.ShippingTask)

	return &Queue{
		server: srv,
		mux:    mux,
	}
}

func (q *Queue) Run() {
	log.Logger.Info("asinq queue server run start")

	if err := q.server.Run(q.mux); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Queue/Run", "section": "server.Run"}).Fatal(err)
		return
	}
}

func (q *Queue) Shutdown() {
	q.server.Shutdown()
}
