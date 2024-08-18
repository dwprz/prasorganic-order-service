package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

type Order interface {
	Create(ctx context.Context, data *dto.TransactionReq) error
	FindManyByUserId(ctx context.Context, data *dto.GetOrdersByCurrentUserReq) (
		*entity.DataWithPaging[[]*entity.OrderWithProducts], error) 
}
