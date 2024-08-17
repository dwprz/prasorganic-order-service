package delivery

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

type Shipper interface {
	ShippingOrder(ctx context.Context, data *entity.OrderWithProducts) (shippingId string, err error)
}
