package repository

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type ProductMock struct {
	mock.Mock
}

func NewProductMock() *ProductMock {
	return &ProductMock{
		Mock: mock.Mock{},
	}
}

func (s *ProductMock) FindByOrderId(ctx context.Context, orderId string) ([]*entity.ProductOrder, error) {
	arguments := s.Mock.Called(ctx, orderId)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).([]*entity.ProductOrder), arguments.Error(1)
}
