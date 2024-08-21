package test

import (
	"context"
	"testing"
	"time"

	grpcclient "github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	restfulclient "github.com/dwprz/prasorganic-order-service/src/core/restful/client"

	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-order-service/src/mock/queue"
	"github.com/dwprz/prasorganic-order-service/src/mock/repository"
	servicemock "github.com/dwprz/prasorganic-order-service/src/mock/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	serviceimpl "github.com/dwprz/prasorganic-order-service/src/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_CreateTransaction$ -v ./src/service/test/ -count=1

type CreateTxTestSuite struct {
	suite.Suite
	txService       service.Transaction
	orderService    *servicemock.OrderMock
	orderRepo       *repository.OrderMock
	midtransRESTful *delivery.MidtransRESTfulMock
}

func (c *CreateTxTestSuite) SetupSuite() {
	productDelivery := delivery.NewProductGrpcMock()
	productConn := new(grpc.ClientConn)

	grpcClient := grpcclient.NewGrpc(productDelivery, productConn)
	c.orderRepo = repository.NewOrderMock()

	c.midtransRESTful = delivery.NewMidtransRESTfulMock()
	shipperDelivery := delivery.NewShipperRESTfulMock()

	restfulClient := restfulclient.NewRestful(c.midtransRESTful, shipperDelivery)

	c.orderService = servicemock.NewOrderMock()
	orderRepo := repository.NewOrderMock()
	productRepo := repository.NewProductMock()

	queueClient := queue.NewClientMock()

	c.txService = serviceimpl.NewTransaction(grpcClient, restfulClient, c.orderService, orderRepo, productRepo, queueClient)
}

func (c *CreateTxTestSuite) Test_Success() {
	req := c.CreateTransactionReq()
	midtransRes := c.CreateMidtransTxRes()

	c.MockMidtransDelivery_Transaction(req, midtransRes, nil)
	c.MockOrderService_Create(req, nil)

	res, err := c.txService.Create(context.Background(), req)
	assert.NoError(c.T(), err)

	assert.NotEmpty(c.T(), res.OrderId)
	assert.Equal(c.T(), midtransRes.Token, res.Token)
	assert.Equal(c.T(), midtransRes.RedirectUrl, res.RedirectUrl)
}

func (c *CreateTxTestSuite) Test_NilProduct() {
	req := c.CreateTransactionReq()
	req.Products = nil

	res, err := c.txService.Create(context.Background(), req)
	assert.Error(c.T(), err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(c.T(), ok)

	assert.Equal(c.T(), validationErr[0].Field(), "Products")
	assert.Nil(c.T(), res)
}

func (c *CreateTxTestSuite) MockMidtransDelivery_Transaction(data *dto.TransactionReq, returnArg1 *dto.MidtransTxRes, returnArg2 error) {

	c.midtransRESTful.Mock.On("Transaction", mock.Anything, mock.MatchedBy(func(req *dto.TransactionReq) bool {
		return data.Order.UserId == req.Order.UserId && req.Order.OrderId != ""
	})).Return(returnArg1, returnArg2)
}

func (c *CreateTxTestSuite) MockOrderService_Create(data *dto.TransactionReq, returnArg1 error) {

	c.orderService.Mock.On("Create", mock.Anything, mock.MatchedBy(func(req *dto.TransactionReq) bool {
		return data.Order.UserId == req.Order.UserId && req.Order.OrderId != "" && req.Order.SnapToken != "" && req.Order.SnapRedirectURL != ""
	})).Return(returnArg1)
}

func (c *CreateTxTestSuite) CreateMidtransTxRes() *dto.MidtransTxRes {

	return &dto.MidtransTxRes{
		Token:       "wjdsajdhweui38wsj2e8uq",
		RedirectUrl: "https://payment.gateway/redirect",
	}
}

func (c *CreateTxTestSuite) CreateTransactionReq() *dto.TransactionReq {
	return &dto.TransactionReq{
		Order: &entity.Order{
			GrossAmount:   500000,
			Status:        entity.PENDING_PAYMENT,
			ShippingId:    "ship-456",
			Courier:       "JNE",
			RateId:        1,
			RateName:      "Express",
			RateType:      "Overnight",
			COD:           false,
			UseInsurance:  true,
			PackageType:   2,
			PaymentMethod: "Credit Card",
			UserId:        "hyfca5Sq7nQcaY6ACksXP",
			Email:         "user@example.com",
			Buyer:         "John Doe",
			Height:        10,
			Length:        20,
			Width:         15,
			Weight:        2.5,
			AddressOwner:  "John Doe",
			Street:        "123 Main St",
			AreaId:        1234,
			Area:          "Central",
			Lat:           "-6.200000",
			Lng:           "106.816666",
			Suburb:        "Jakarta Pusat",
			City:          "Jakarta",
			Province:      "DKI Jakarta",
			WhatsApp:      "081234567890",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
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

func TestService_CreateTransaction(t *testing.T) {
	suite.Run(t, new(CreateTxTestSuite))
}
