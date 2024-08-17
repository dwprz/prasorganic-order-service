package delivery

import (
	"context"
	"log"

	"github.com/dwprz/prasorganic-order-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/interface/delivery"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGrpcImpl struct {
	client pb.ProductServiceClient
}

func NewProductGrpc(unaryRequest *interceptor.UnaryRequest) (delivery.ProductGrpc, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryRequest.AddBasicAuth),
	)

	conn, err := grpc.NewClient(config.Conf.ApiGateway.BaseUrl, opts...)
	if err != nil {
		log.Fatalf("new otp grpc client: %v", err.Error())
	}

	client := pb.NewProductServiceClient(conn)

	return &ProductGrpcImpl{
		client: client,
	}, conn
}

func (p *ProductGrpcImpl) ReduceStocks(ctx context.Context, data []*entity.ProductOrder) error {
	var req []*pb.ProductOrder
	if err := copier.Copy(&req, data); err != nil {
		return err
	}

	_, err := cbreaker.ProductGrpc.Execute(func() (any, error) {
		_, err := p.client.ReduceStocks(ctx, &pb.ReduceStocksReq{
			Data: req,
		})

		return nil, err
	})

	return err
}

func (p *ProductGrpcImpl) RollbackStocks(ctx context.Context, data []*entity.ProductOrder) error {
	var req []*pb.ProductOrder
	if err := copier.Copy(&req, data); err != nil {
		return err
	}

	_, err := cbreaker.ProductGrpc.Execute(func() (any, error) {
		_, err := p.client.RollbackStocks(ctx, &pb.RollbackStocksReq{
			Data: req,
		})

		return nil, err
	})

	return err
}
