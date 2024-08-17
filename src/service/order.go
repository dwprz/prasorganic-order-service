package service

import (
	"context"

	v "github.com/dwprz/prasorganic-order-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
)

type OrderImpl struct {
	orderRepo repository.Order
}

func NewOrder(or repository.Order) service.Order {
	return &OrderImpl{
		orderRepo: or,
	}
}

func (o *OrderImpl) Create(ctx context.Context, data *dto.TransactionReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	for _, product := range data.Products {
		product.OrderId = data.Order.OrderId
	}

	err := o.orderRepo.Create(ctx, data)
	return err
}

