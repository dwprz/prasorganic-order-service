package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-order-service/src/core/grpc"
	"github.com/dwprz/prasorganic-order-service/src/core/restful"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/repository"
	"github.com/dwprz/prasorganic-order-service/src/service"
)

func handleCloseApp(closeCH chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		close(closeCH)
	}()
}

func main() {
	closeCH := make(chan struct{})
	handleCloseApp(closeCH)

	postgresDB := database.NewPostgres()

	restfulClient := restful.InitClient()
	grpcClient := grpc.InitClient()

	orderRepo := repository.NewOrder(postgresDB, grpcClient)

	orderService := service.NewOrder(orderRepo)
	txService := service.NewTransaction(restfulClient, orderService)

	restfulServer := restful.InitServer(txService)
	defer restfulServer.Stop()

	go restfulServer.Run()

	<-closeCH
}
