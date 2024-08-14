package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
)

type Transaction interface {
	Create(ctx context.Context, data *dto.TransactionReq) (*dto.TransactionRes, error)
}
