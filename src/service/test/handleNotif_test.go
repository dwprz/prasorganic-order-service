package test

import (
	"context"
	"testing"

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

// go test -p=1 -v ./src/service/test/... -count=1
// go test -run ^TestService_HandleNotif$ -v ./src/service/test/ -count=1

type HandleNotifTestSuite struct {
	suite.Suite
	txService   service.Transaction
	productGrpc *delivery.ProductGrpcMock
	orderRepo   *repository.OrderMock
	productRepo *repository.ProductMock
}

func (h *HandleNotifTestSuite) SetupSuite() {
	h.productGrpc = delivery.NewProductGrpcMock()
	productConn := new(grpc.ClientConn)

	grpcClient := grpcclient.NewGrpc(h.productGrpc, productConn)

	midtransRESTful := delivery.NewMidtransRESTfulMock()
	shipperRESTful := delivery.NewShipperRESTfulMock()

	restfulClient := restfulclient.NewRestful(midtransRESTful, shipperRESTful)

	orderService := servicemock.NewOrderMock()
	h.orderRepo = repository.NewOrderMock()
	h.productRepo = repository.NewProductMock()

	queueClient := queue.NewClientMock()

	h.txService = serviceimpl.NewTransaction(grpcClient, restfulClient, orderService, h.orderRepo, h.productRepo, queueClient)
}

func (h *HandleNotifTestSuite) Test_PAID() {
	req := h.CreateTransaction()
	req.OrderId = "hyfa_5Sq7nQcaY6ACksXP"
	req.TransactionStatus = "capture"
	req.FraudStatus = "accept"

	h.MockOrderRepo_UpdateById(req.OrderId, entity.PAID, req.PaymentType, nil)

	err := h.txService.HandleNotif(context.Background(), req)
	assert.NoError(h.T(), err)
}

func (h *HandleNotifTestSuite) Test_CANCEL() {
	req := h.CreateTransaction()
	req.OrderId = "hija_5Sq7nQcaY6ACksXP"
	req.TransactionStatus = "cancel"

	products := h.CreateProductOrders(req.OrderId)
	h.MockProductRepo_FindByOrderId(req.OrderId, products, nil)

	h.MockProductGrpc_Rollbacstocks(products, nil)
	h.MockOrderRepo_UpdateById(req.OrderId, entity.CANCELLED, "", nil)

	err := h.txService.HandleNotif(context.Background(), req)
	assert.NoError(h.T(), err)
}

func (h *HandleNotifTestSuite) Test_FAILED() {
	req := h.CreateTransaction()
	req.OrderId = "klma_5Sq7nQcaY6ACksXP"
	req.TransactionStatus = "expire"

	products := h.CreateProductOrders(req.OrderId)
	h.MockProductRepo_FindByOrderId(req.OrderId, products, nil)

	h.MockProductGrpc_Rollbacstocks(products, nil)
	h.MockOrderRepo_UpdateById(req.OrderId, entity.FAILED, "", nil)

	err := h.txService.HandleNotif(context.Background(), req)
	assert.NoError(h.T(), err)
}

func (h *HandleNotifTestSuite) MockProductGrpc_Rollbacstocks(data []*entity.ProductOrder, returnArg1 error) {

	h.productGrpc.Mock.On("RollbackStocks", mock.Anything, data).Return(returnArg1)
}

func (h *HandleNotifTestSuite) MockOrderRepo_UpdateById(orderId string, status entity.OrderStatus,
	paymentMethod string, returnArg1 error) {

	h.orderRepo.Mock.On("UpdateById", mock.Anything, &entity.Order{
		OrderId:       orderId,
		Status:        status,
		PaymentMethod: paymentMethod,
	}).Return(returnArg1)
}

func (h *HandleNotifTestSuite) MockProductRepo_FindByOrderId(orderId string, returnArg1 []*entity.ProductOrder, returnArg2 error) {

	h.productRepo.Mock.On("FindByOrderId", mock.Anything, orderId).Return(returnArg1, returnArg2)
}

func (h *HandleNotifTestSuite) CreateMidtransTxRes() *dto.MidtransTxRes {

	return &dto.MidtransTxRes{
		Token:       "wjdsajdhweui38wsj2e8uq",
		RedirectUrl: "https://payment.gateway/redirect",
	}
}

func (h *HandleNotifTestSuite) CreateTransaction() *entity.Transaction {
	return &entity.Transaction{
		TransactionTime:        "2024-08-21 12:34:56",
		TransactionId:          "txn-001",
		StatusMessage:          "Transaction successful",
		StatusCode:             "200",
		SignatureKey:           "abcdef123456",
		PaymentType:            "credit_card",
		OrderId:                "hyfa_5Sq7nQcaY6ACksXP",
		MerchantId:             "merchant-123",
		MaskedCard:             "4111-XXXX-XXXX-1111",
		GrossAmount:            "100.00",
		FraudStatus:            "accept",
		Eci:                    "05",
		Currency:               "IDR",
		ChannelResponseMessage: "Approved",
		ChannelResponseCode:    "00",
		CardType:               "visa",
		Bank:                   "BCA",
		ApprovalCode:           "123456",
	}
}

func (h *HandleNotifTestSuite) CreateProductOrders(orderId string) []*entity.ProductOrder {
	return []*entity.ProductOrder{
		{
			OrderId:     orderId,
			ProductId:   1,
			ProductName: "Product A",
			Image:       "https://example.com/image-a.jpg",
			Quantity:    2,
			Price:       150000,
		},
		{
			OrderId:     orderId,
			ProductId:   2,
			ProductName: "Product B",
			Image:       "https://example.com/image-b.jpg",
			Quantity:    1,
			Price:       200000,
		},
	}
}

func TestService_HandleNotif(t *testing.T) {
	suite.Run(t, new(HandleNotifTestSuite))
}
