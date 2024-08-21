package service

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type Transaction interface {
	Create(ctx context.Context, data *dto.TransactionReq) (*dto.TransactionRes, error)
	HandleNotif(ctx context.Context, data *entity.Transaction) error
}

type TransactionMock struct {
	mock.Mock
}

func NewTransactionMock() *TransactionMock {
	return &TransactionMock{
		Mock: mock.Mock{},
	}
}

func (s *TransactionMock) Create(ctx context.Context, data *dto.TransactionReq) (*dto.TransactionRes, error) {
	arguments := s.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.TransactionRes), arguments.Error(1)
}

func (s *TransactionMock) HandleNotif(ctx context.Context, data *entity.Transaction) error {
	arguments := s.Mock.Called(ctx, data)

	return arguments.Error(0)
}
