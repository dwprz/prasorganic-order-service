package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	v "github.com/dwprz/prasorganic-order-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
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

func (o *OrderImpl) FindManyByUserId(ctx context.Context, data *dto.GetOrdersByCurrentUserReq) (
	*entity.DataWithPaging[[]*entity.OrderWithProducts], error) {

	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	limit, offset := helper.CreateLimitAndOffset(data.Page)

	res, err := o.orderRepo.FindManyByUserId(ctx, data.UserId, limit, offset)
	if err != nil {
		return nil, err
	}

	return helper.FormatPagedData(res.Orders, res.TotalOrders, data.Page, limit), nil
}
