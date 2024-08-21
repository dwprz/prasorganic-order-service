package delivery

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/stretchr/testify/mock"
)

type MidtransRESTfulMock struct {
	mock.Mock
}

func NewMidtransRESTfulMock() *MidtransRESTfulMock {
	return &MidtransRESTfulMock{
		Mock: mock.Mock{},
	}
}

func (m *MidtransRESTfulMock) Transaction(ctx context.Context, data *dto.TransactionReq) (*dto.MidtransTxRes, error) {
	arguments := m.Mock.Called(ctx, data)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*dto.MidtransTxRes), arguments.Error(1)
}
