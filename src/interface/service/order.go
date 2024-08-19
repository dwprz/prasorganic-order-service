package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

type Order interface {
	Create(ctx context.Context, data *dto.TransactionReq) error
	FindManyByUserId(ctx context.Context, data *dto.GetOrdersByCurrentUserReq) (*entity.DataWithPaging[[]*entity.OrderWithProducts], error)
	FindMany(ctx context.Context, data *dto.GetOrdersReq) (*entity.DataWithPaging[[]*entity.OrderWithProducts], error)
	Cancel(ctx context.Context, data *dto.CancelOrderReq) error
	UpdateStatus(ctx context.Context, data *dto.UpdateStatusReq) error
	AddShippingId(ctx context.Context, data *dto.AddShippingIdReq) error
}
