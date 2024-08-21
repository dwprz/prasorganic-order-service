package test

import (
	"context"
	"testing"
	"time"

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
// go test -run ^TestService_Cancel$ -v ./src/service/test/ -count=1

type CancelTestSuite struct {
	suite.Suite
	orderService service.Order
	orderRepo    *repository.OrderMock
}

func (c *CancelTestSuite) SetupSuite() {
	c.orderRepo = repository.NewOrderMock()

	c.orderService = serviceimpl.NewOrder(c.orderRepo)
}

func (c *CancelTestSuite) Test_Success() {
	userId := "syfhi5Sq7nQcaY6ACksXP"
	orderId := "hyfde5Sq7nQcaY6ACksXP"

	req := c.CreateCancelOrdersReq(userId, orderId)

	order := c.CreateOrderWithProducts()
	order.Order.UserId = userId
	order.Order.OrderId = orderId
	order.Order.ShippingId = ""

	c.orderRepo.Mock.On("FindById", mock.Anything, orderId).Return(order, nil)
	c.MockOrderRepo_UpdateById(orderId, nil)

	err := c.orderService.Cancel(context.Background(), req)
	assert.NoError(c.T(), err)
}

func (c *CancelTestSuite) Test_InvalidUserId() {
	userId := "invalid-user-id"
	orderId := "defde5Sq7nQcaY6ACksXP"

	req := c.CreateCancelOrdersReq(userId, orderId)

	err := c.orderService.Cancel(context.Background(), req)
	assert.Error(c.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(c.T(), ok)

	assert.Equal(c.T(), validationErr[0].Field(), "UserId")
}

func (c *CancelTestSuite) MockOrderRepo_UpdateById(orderId string, returnArg1 error) {
	c.orderRepo.Mock.On("UpdateById", mock.Anything, &entity.Order{
		OrderId: orderId,
		Status:  entity.CANCELLED,
	}).Return(returnArg1)
}

func (c *CancelTestSuite) CreateCancelOrdersReq(userId string, orderId string) *dto.CancelOrderReq {
	return &dto.CancelOrderReq{
		UserId:  userId,
		OrderId: orderId,
	}
}

func (c *CancelTestSuite) CreateOrderWithProducts() *entity.OrderWithProducts {
	return &entity.OrderWithProducts{
		Order: &entity.Order{
			OrderId:         "hyfa_5Sq7nQcaY6ACksXP",
			GrossAmount:     500000,
			Status:          entity.PENDING_PAYMENT,
			ShippingId:      "ship-456",
			Courier:         "JNE",
			RateId:          1,
			RateName:        "Express",
			RateType:        "Overnight",
			COD:             false,
			UseInsurance:    true,
			PackageType:     2,
			PaymentMethod:   "Credit Card",
			SnapToken:       "snap-token-789",
			SnapRedirectURL: "https://payment.gateway/redirect",
			UserId:          "hyfa_5Sq7nQcaY6ACkabc",
			Email:           "user@example.com",
			Buyer:           "John Doe",
			Height:          10,
			Length:          20,
			Width:           15,
			Weight:          2.5,
			AddressOwner:    "John Doe",
			Street:          "123 Main St",
			AreaId:          1234,
			Area:            "Central",
			Lat:             "-6.200000",
			Lng:             "106.816666",
			Suburb:          "Jakarta Pusat",
			City:            "Jakarta",
			Province:        "DKI Jakarta",
			WhatsApp:        "081234567890",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		Products: []*entity.ProductOrder{
			{
				OrderId:     "order-123",
				ProductId:   1,
				ProductName: "Product A",
				Image:       "https://example.com/image-a.jpg",
				Quantity:    2,
				Price:       150000,
			},
			{
				OrderId:     "order-123",
				ProductId:   2,
				ProductName: "Product B",
				Image:       "https://example.com/image-b.jpg",
				Quantity:    1,
				Price:       200000,
			},
		},
	}
}

func TestService_Cancel(t *testing.T) {
	suite.Run(t, new(CancelTestSuite))
}
