package service

import (
	"context"
	"strings"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
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

func (o *OrderImpl) FindMany(ctx context.Context, data *dto.GetOrdersReq) (
	*entity.DataWithPaging[[]*entity.OrderWithProducts], error) {

	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	limit, offset := helper.CreateLimitAndOffset(data.Page)

	var res *dto.OrdersWithCountRes
	var err error

	switch {
	case data.Status != "":
		data.Status = strings.ToUpper(data.Status)
		res, err = o.orderRepo.FindManyByStatus(ctx, data.Status, limit, offset)
	default:
		res, err = o.orderRepo.FindMany(ctx, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	return helper.FormatPagedData(res.Orders, res.TotalOrders, data.Page, limit), nil
}

func (o *OrderImpl) Cancel(ctx context.Context, data *dto.CancelOrderReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	res, err := o.orderRepo.FindById(ctx, data.OrderId)
	if err != nil {
		return err
	}

	if res.Order.UserId != data.UserId {
		return &errors.Response{HttpCode: 403, Message: "you do not have permission to cancel this order"}
	}

	if res.Order.ShippingId != "" {
		return &errors.Response{HttpCode: 400, Message: "sorry, your order has been processed"}
	}

	err = o.orderRepo.UpdateById(ctx, &entity.Order{
		OrderId: data.OrderId,
		Status:  entity.CANCELLED,
	})

	return err
}

func (o *OrderImpl) UpdateStatus(ctx context.Context, data *dto.UpdateStatusReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	err := o.orderRepo.UpdateById(ctx, &entity.Order{
		OrderId: data.OrderId,
		Status:  entity.OrderStatus(data.Status),
	})

	return err
}

func (o *OrderImpl) AddShippingId(ctx context.Context, data *dto.AddShippingIdReq) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	err := o.orderRepo.UpdateById(ctx, &entity.Order{
		OrderId:    data.OrderId,
		ShippingId: data.ShippingId,
	})

	return err
}
