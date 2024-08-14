package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
)

type Order interface {
	Create(ctx context.Context, data *dto.TransactionReq) error
}
