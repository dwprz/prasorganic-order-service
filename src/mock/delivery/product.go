package delivery

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/stretchr/testify/mock"
)

type ProductGrpcMock struct {
	mock.Mock
}

func NewProductGrpcMock() *ProductGrpcMock {
	return &ProductGrpcMock{
		Mock: mock.Mock{},
	}
}

func (p *ProductGrpcMock) ReduceStocks(ctx context.Context, data []*entity.ProductOrder) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}

func (p *ProductGrpcMock) RollbackStocks(ctx context.Context, data []*entity.ProductOrder) error {
	arguments := p.Mock.Called(ctx, data)

	return arguments.Error(0)
}
