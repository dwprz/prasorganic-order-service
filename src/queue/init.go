package queue

import (
	restfulclient "github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-order-service/src/interface/queue"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/queue/client"
	"github.com/dwprz/prasorganic-order-service/src/queue/handler"
	"github.com/dwprz/prasorganic-order-service/src/queue/server"
)

func InitServer(rc *restfulclient.Restful, or repository.Order) *server.Queue {
	orderHandler := handler.NewOrderQueue(rc, or)
	queueServer := server.NewQueue(orderHandler)

	return queueServer
}

func InitClient() queue.Client {
	queueClient := client.NewQueue()

	return queueClient
}
