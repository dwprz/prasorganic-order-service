package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type OrderMock struct {
	mock.Mock
}

func NewOrderMock() *OrderMock {
	return &OrderMock{
		Mock: mock.Mock{},
	}
}

func (s *OrderMock) Create(ctx context.Context, data *dto.TransactionReq) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (s *OrderMock) FindManyByUserId(ctx context.Context, data *dto.GetOrdersByCurrentUserReq) (*entity.DataWithPaging[[]*entity.OrderWithProducts], error) {
	arguments := s.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.DataWithPaging[[]*entity.OrderWithProducts]), arguments.Error(1)
}

func (s *OrderMock) FindMany(ctx context.Context, data *dto.GetOrdersReq) (*entity.DataWithPaging[[]*entity.OrderWithProducts], error) {
	arguments := s.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.DataWithPaging[[]*entity.OrderWithProducts]), arguments.Error(1)
}

func (s *OrderMock) Cancel(ctx context.Context, data *dto.CancelOrderReq) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (s *OrderMock) UpdateStatus(ctx context.Context, data *dto.UpdateStatusReq) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (s *OrderMock) AddShippingId(ctx context.Context, data *dto.AddShippingIdReq) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}