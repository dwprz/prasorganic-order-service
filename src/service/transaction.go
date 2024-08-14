package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	v "github.com/dwprz/prasorganic-order-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type TransactionImpl struct {
	restfulClient *client.Restful
	orderService  service.Order
}

func NewTransaction(rc *client.Restful, os service.Order) service.Transaction {
	return &TransactionImpl{
		restfulClient: rc,
		orderService:  os,
	}
}

func (t *TransactionImpl) Create(ctx context.Context, data *dto.TransactionReq) (*dto.TransactionRes, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	orderId, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	data.Order.OrderId = orderId
	txRes, err := t.restfulClient.Midtrans.Transaction(ctx, data)
	if err != nil {
		return nil, err
	}

	data.Order.SnapToken = txRes.Token
	data.Order.SnapRedirectURL = txRes.RedirectUrl

	helper.LogJSON(data)

	if err := t.orderService.Create(ctx, data); err != nil {
		return nil, err
	}

	return &dto.TransactionRes{
		OrderId:     orderId,
		Token:       txRes.Token,
		RedirectUrl: txRes.RedirectUrl,
	}, nil
}
