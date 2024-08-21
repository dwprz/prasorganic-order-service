package restful

import (
	"github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/delivery"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
)

func InitServer(ts service.Transaction, os service.Order) *server.Restful {
	orderHandler := handler.NewOrderRESTful(ts, os)
	middleware := middleware.New()

	restfulServer := server.NewRestful(orderHandler, middleware)
	return restfulServer
}

func InitClient() *client.Restful {
	midtransDelivery := delivery.NewMidtransRESTful()
	shipperDelivery := delivery.NewShipperRESTful()
	restfulClient := client.NewRestful(midtransDelivery, shipperDelivery)

	return restfulClient
}
