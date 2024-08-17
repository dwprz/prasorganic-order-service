package restful

import (
	"github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
)

func InitServer(ts service.Transaction) *server.Restful {
	txHandler := handler.NewTransaction(ts)
	middleware := middleware.New()

	restfulServer := server.NewRestful(txHandler, middleware)
	return restfulServer
}

func InitClient() *client.Restful {
	midtransDelivery := delivery.NewMidtrans()
	shipperDelivery := delivery.NewShipper()
	restfulClient := client.NewRestful(midtransDelivery, shipperDelivery)

	return restfulClient
}
