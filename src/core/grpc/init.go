package grpc

import (
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/delivery"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/interceptor"
)

func InitClient() *client.Grpc {
	unaryRequestInterceptor := interceptor.NewUnaryRequest()
	productDelivery, productConn := delivery.NewProductGrpc(unaryRequestInterceptor)

	grpcClient := client.NewGrpc(productDelivery, productConn)
	return grpcClient
}
