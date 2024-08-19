package handler

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	pb "github.com/dwprz/prasorganic-proto/protogen/order"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OtpGrpcImpl struct {
	orderService service.Order
	pb.UnimplementedOrderServiceServer
}

func NewOrderGrpc(os service.Order) pb.OrderServiceServer {
	return &OtpGrpcImpl{
		orderService: os,
	}
}

func (a *OtpGrpcImpl) AddShippingId(ctx context.Context, data *pb.AddShippingIdReq) (*emptypb.Empty, error) {
	err := a.orderService.AddShippingId(ctx, &dto.AddShippingIdReq{
		OrderId:    data.OrderId,
		ShippingId: data.ShippingId,
	})

	return nil, err
}

func (a *OtpGrpcImpl) UpdateStatus(ctx context.Context, data *pb.UpdateStatusReq) (*emptypb.Empty, error) {
	err := a.orderService.UpdateStatus(ctx, &dto.UpdateStatusReq{
		OrderId: data.OrderId,
		Status:  data.Status,
	})

	return nil, err
}
