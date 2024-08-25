package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	grpcclient "github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	restfulclient "github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	v "github.com/dwprz/prasorganic-order-service/src/infrastructure/validator"
	"github.com/dwprz/prasorganic-order-service/src/interface/queue"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type TransactionImpl struct {
	grpcClient    *grpcclient.Grpc
	restfulClient *restfulclient.Restful
	orderService  service.Order
	orderRepo     repository.Order
	productRepo   repository.Product
	queueClient   queue.Client
}

func NewTransaction(gc *grpcclient.Grpc, rc *restfulclient.Restful, os service.Order, or repository.Order,
	pr repository.Product, qc queue.Client) service.Transaction {
	return &TransactionImpl{
		grpcClient:    gc,
		restfulClient: rc,
		orderService:  os,
		orderRepo:     or,
		productRepo:   pr,
		queueClient:   qc,
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

	data.Order.OrderId = strings.ToUpper(orderId)
	txRes, err := t.restfulClient.Midtrans.Transaction(ctx, data)
	if err != nil {
		return nil, err
	}

	data.Order.SnapToken = txRes.Token
	data.Order.SnapRedirectURL = txRes.RedirectUrl

	if err := t.orderService.Create(ctx, data); err != nil {
		return nil, err
	}

	return &dto.TransactionRes{
		OrderId:     orderId,
		Token:       txRes.Token,
		RedirectUrl: txRes.RedirectUrl,
	}, nil
}

func (t *TransactionImpl) HandleNotif(ctx context.Context, data *entity.Transaction) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	if (data.TransactionStatus == "capture" && data.FraudStatus == "accept") ||
		data.TransactionStatus == "settlement" {

		err := t.orderRepo.UpdateById(ctx, &entity.Order{
			OrderId:       data.OrderId,
			Status:        entity.PAID,
			PaymentMethod: data.PaymentType,
		})

		if err != nil {
			return err
		}

		payload, err := json.Marshal(map[string]string{"order_id": data.OrderId})
		if err != nil {
			return err
		}

		t.queueClient.Create("orders:shipping", "orders", payload, time.Duration(1*time.Hour))
	}

	if data.TransactionStatus == "cancel" {
		products, err := t.productRepo.FindByOrderId(ctx, data.OrderId)
		if err != nil {
			return err
		}

		if err := t.grpcClient.Product.RollbackStocks(ctx, products); err != nil {
			return err
		}

		err = t.orderRepo.UpdateById(ctx, &entity.Order{
			OrderId: data.OrderId,
			Status:  entity.CANCELLED,
		})

		return err
	}

	if data.TransactionStatus == "deny" || data.TransactionStatus == "expire" {
		products, err := t.productRepo.FindByOrderId(ctx, data.OrderId)
		if err != nil {
			return err
		}

		if err := t.grpcClient.Product.RollbackStocks(ctx, products); err != nil {
			return err
		}

		err = t.orderRepo.UpdateById(ctx, &entity.Order{
			OrderId: data.OrderId,
			Status:  entity.FAILED,
		})

		return err
	}

	return nil
}
