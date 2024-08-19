package grpc

import (
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/delivery"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/handler"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/server"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
)

func InitClient() *client.Grpc {
	unaryRequestInterceptor := interceptor.NewUnaryRequest()
	productDelivery, productConn := delivery.NewProductGrpc(unaryRequestInterceptor)

	grpcClient := client.NewGrpc(productDelivery, productConn)
	return grpcClient
}

func InitServer(os service.Order) *server.Grpc {
	orderHandler := handler.NewOrderGrpc(os)
	unaryResponseInterceptor := interceptor.NewUnaryResponse()

	grpcServer := server.NewGrpc(orderHandler, unaryResponseInterceptor)
	return grpcServer
}
