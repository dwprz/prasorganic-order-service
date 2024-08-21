package delivery

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type ShipperRESTfulMock struct {
	mock.Mock
}

func NewShipperRESTfulMock() *ShipperRESTfulMock {
	return &ShipperRESTfulMock{
		Mock: mock.Mock{},
	}
}

func (s *ShipperRESTfulMock) ShippingOrder(ctx context.Context, data *entity.OrderWithProducts) (shippingId string, err error) {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Get(0).(string), arguments.Error(1)
}