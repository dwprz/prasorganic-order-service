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
// go test -run ^TestService_FindMany$ -v ./src/service/test/ -count=1

type FindManyTestSuite struct {
	suite.Suite
	orderService service.Order
	orderRepo    *repository.OrderMock
}

func (f *FindManyTestSuite) SetupSuite() {
	f.orderRepo = repository.NewOrderMock()

	f.orderService = serviceimpl.NewOrder(f.orderRepo)
}

func (f *FindManyTestSuite) Test_Success() {
	req := f.CreateGetOrdersReq("", 1)

	orderWithCountRes := f.CreateOrderWithCountRes()
	f.orderRepo.Mock.On("FindMany", mock.Anything, 20, 0).Return(orderWithCountRes, nil)

	res, err := f.orderService.FindMany(context.Background(), req)
	assert.NoError(f.T(), err)

	assert.NotEmpty(f.T(), res)
}

func (f *FindManyTestSuite) Test_LimitPage() {
	req := f.CreateGetOrdersReq("PENDING_PAYMENT", 1000)

	res, err := f.orderService.FindMany(context.Background(), req)
	assert.Error(f.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(f.T(), ok)

	assert.Equal(f.T(), validationErr[0].Field(), "Page")
	assert.Nil(f.T(), res)
}

func (f *FindManyTestSuite) CreateGetOrdersReq(status string, page int) *dto.GetOrdersReq {
	return &dto.GetOrdersReq{
		Status: status,
		Page:   page,
	}
}

func (f *FindManyTestSuite) CreateOrderWithProducts() *entity.OrderWithProducts {
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

func (f *FindManyTestSuite) CreateOrderWithCountRes() *dto.OrdersWithCountRes {
	order := f.CreateOrderWithProducts()
	orders := []*entity.OrderWithProducts{order}

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: len(orders),
	}
}

func TestService_FindMany(t *testing.T) {
	suite.Run(t, new(FindManyTestSuite))
}
