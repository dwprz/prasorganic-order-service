package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/mock/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-order-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_AddShippingId$ -v ./src/service/test/ -count=1

type AddShippingIdTestSuite struct {
	suite.Suite
	orderService service.Order
	orderRepo    *repository.OrderMock
}

func (a *AddShippingIdTestSuite) SetupSuite() {
	a.orderRepo = repository.NewOrderMock()

	a.orderService = serviceimpl.NewOrder(a.orderRepo)
}

func (a *AddShippingIdTestSuite) Test_Success() {
	orderId := "hyfde5Sq7nQcaY6ACksXP"
	shippingId := "syfhi5Sq7nQcaY6ACksXP"

	req := a.CreateAddShippingIdReq(orderId, shippingId)
	a.MockOrderRepo_UpdateById(orderId, shippingId, nil)

	err := a.orderService.AddShippingId(context.Background(), req)
	assert.NoError(a.T(), err)
}

func (a *AddShippingIdTestSuite) Test_InvalidOrderId() {
	orderId := "invalid-order-id"
	shippingId := "syfhi5Sq7nQcaY6ACksXP"

	req := a.CreateAddShippingIdReq(orderId, shippingId)

	err := a.orderService.AddShippingId(context.Background(), req)
	assert.Error(a.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(a.T(), ok)

	assert.Equal(a.T(), validationErr[0].Field(), "OrderId")
}

func (a *AddShippingIdTestSuite) MockOrderRepo_UpdateById(orderId string, shippingId string, returnArg1 error) {
	a.orderRepo.Mock.On("UpdateById", mock.Anything, &entity.Order{
		OrderId:    orderId,
		ShippingId: shippingId,
	}).Return(returnArg1)
}

func (a *AddShippingIdTestSuite) CreateAddShippingIdReq(orderId string, shippingId string) *dto.AddShippingIdReq {
	return &dto.AddShippingIdReq{
		OrderId:    orderId,
		ShippingId: shippingId,
	}
}

func TestService_AddShippingId(t *testing.T) {
	suite.Run(t, new(AddShippingIdTestSuite))
}
