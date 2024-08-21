package repository

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

func (s *OrderMock) FindById(ctx context.Context, orderId string) (*entity.OrderWithProducts, error) {
	arguments := s.Mock.Called(ctx, orderId)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*entity.OrderWithProducts), arguments.Error(1)
}

func (s *OrderMock) FindMany(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {
	arguments := s.Mock.Called(ctx, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.OrdersWithCountRes), arguments.Error(1)
}

func (s *OrderMock) FindManyByUserId(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {
	arguments := s.Mock.Called(ctx, userId, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.OrdersWithCountRes), arguments.Error(1)
}

func (s *OrderMock) FindManyByStatus(ctx context.Context, status string, limit, offset int) (*dto.OrdersWithCountRes, error) {
	arguments := s.Mock.Called(ctx, status, limit, offset)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.OrdersWithCountRes), arguments.Error(1)
}

func (s *OrderMock) UpdateById(ctx context.Context, data *entity.Order) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}
