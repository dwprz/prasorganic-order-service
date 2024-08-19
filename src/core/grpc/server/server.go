package server

import (
	"fmt"
	"net"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	pb "github.com/dwprz/prasorganic-proto/protogen/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Grpc struct {
	port                     string
	server                   *grpc.Server
	orderHandler             pb.OrderServiceServer
	unaryResponseInterceptor *interceptor.UnaryResponse
}

// this main grpc server
func NewGrpc(orderHandler pb.OrderServiceServer, uri *interceptor.UnaryResponse) *Grpc {
	port := config.Conf.CurrentApp.GrpcPort

	return &Grpc{
		port:                     port,
		orderHandler:             orderHandler,
		unaryResponseInterceptor: uri,
	}
}

func (s *Grpc) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "net.Listen"}).Fatal(err)
	}

	log.Logger.Infof("grpc run in port: %s", s.port)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.unaryResponseInterceptor.Recovery,
			s.unaryResponseInterceptor.Error,
		))

	s.server = grpcServer

	pb.RegisterOrderServiceServer(grpcServer, s.orderHandler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "server.Grpc/Run", "section": "grpcServer.Serve"}).Fatal(err)
	}
}

func (s *Grpc) Stop() {
	s.server.Stop()
}
