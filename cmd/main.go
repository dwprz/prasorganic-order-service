package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwprz/prasorganic-order-service/src/core/broker"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc"
	"github.com/dwprz/prasorganic-order-service/src/core/restful"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/queue"
	"github.com/dwprz/prasorganic-order-service/src/repository"
	"github.com/dwprz/prasorganic-order-service/src/service"
)

func handleCloseApp(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	postgresDB := database.NewPostgres()

	restfulClient := restful.InitClient()
	grpcClient := grpc.InitClient()
	defer grpcClient.Close()

	orderRepo := repository.NewOrder(postgresDB, grpcClient)
	productRepo := repository.NewProduct(postgresDB)

	queueServer := queue.InitServer(restfulClient, orderRepo)
	defer queueServer.Shutdown()

	go queueServer.Run()

	queueCLient := queue.InitClient()
	defer queueCLient.Close()

	orderService := service.NewOrder(orderRepo)
	txService := service.NewTransaction(grpcClient, restfulClient, orderService, orderRepo, productRepo, queueCLient)

	kafkaConsumer := broker.InitKafkaConsumer(txService)
	defer kafkaConsumer.Close()

	go kafkaConsumer.Consume(ctx)

	restfulServer := restful.InitServer(txService, orderService)
	defer restfulServer.Stop()

	go restfulServer.Run()

	grpcServer := grpc.InitServer(orderService)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-ctx.Done()
}
