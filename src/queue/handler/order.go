package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/client"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type OrderQueue struct {
	restfulClient *client.Restful
	orderRepo     repository.Order
}

func NewOrderQueue(rc *client.Restful, or repository.Order) *OrderQueue {
	return &OrderQueue{
		restfulClient: rc,
		orderRepo:     or,
	}
}

func (o *OrderQueue) ShippingTask(ctx context.Context, t *asynq.Task) error {
	b := t.Payload()

	var data map[string]string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	orderId := data["order_id"]

	ctx, cancel := context.WithTimeout(ctx, time.Duration(30*time.Second))
	defer cancel()

	order, err := o.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return err
	}

	shippingId, err := o.restfulClient.Shipper.ShippingOrder(ctx, order)
	if err != nil {
		return err
	}

	err = o.orderRepo.UpdateById(ctx, &entity.Order{
		OrderId:    orderId,
		Status:     entity.IN_PROGRESS,
		ShippingId: shippingId,
	})

	log.Logger.WithFields(logrus.Fields{"location": "handler.OrderQueue/ShippingTask"}).Infof("run job successfully")
	return err
}
