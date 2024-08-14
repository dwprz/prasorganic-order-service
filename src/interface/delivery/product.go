package delivery

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

type ProductGrpc interface {
	UpdateStock(ctx context.Context, data []*entity.ProductOrder) error
}
