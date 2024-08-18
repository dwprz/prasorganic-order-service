package repository

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

type Order interface {
	Create(ctx context.Context, data *dto.TransactionReq) error
	FindById(ctx context.Context, orderId string) (*entity.OrderWithProducts, error)
	FindMany(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error)
	FindManyByUserId(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error)
	FindManyByStatus(ctx context.Context, status string, limit, offset int) (*dto.OrdersWithCountRes, error)
	UpdateById(ctx context.Context, data *entity.Order) error
}
