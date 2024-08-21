package test

import (
	"context"
	"testing"
	"time"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	repoimpl "github.com/dwprz/prasorganic-order-service/src/repository"
	"github.com/dwprz/prasorganic-order-service/test/util"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// go test -p=1 -v ./src/repository/test/... -count=1
// go test -run ^TestRepository_CreateOrder$ -v ./src/repository/test/ -count=1

type CreateOrderTestSuite struct {
	suite.Suite
	orderRepo     repository.Order
	productGrpc   *delivery.ProductGrpcMock
	postgresDB    *gorm.DB
	orderTestUtil *util.OrderTest
}

func (c *CreateOrderTestSuite) SetupSuite() {
	c.postgresDB = database.NewPostgres()
	c.productGrpc = delivery.NewProductGrpcMock()
	productConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(c.productGrpc, productConn)

	c.orderRepo = repoimpl.NewOrder(c.postgresDB, grpcClient)

	c.orderTestUtil = util.NewOrderTest(c.postgresDB)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	c.orderTestUtil.Delete()

	sqlDB, err := c.postgresDB.DB()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.CreateOrderTestSuite/TearDownSuite", "section": "postgresDB.DB"}).Fatal(err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.CreateOrderTestSuite/TearDownSuite", "section": "sqlDB.Close"}).Fatal(err)
	}
}

func (c *CreateOrderTestSuite) Test_Success() {
	req := c.CreateTransactionReq()
	c.MockProductGrpc_Rollbacstocks(req.Products, nil)

	err := c.orderRepo.Create(context.Background(), req)
	assert.NoError(c.T(), err)
}

func (c *CreateOrderTestSuite) MockProductGrpc_Rollbacstocks(data []*entity.ProductOrder, returnArg1 error) {

	c.productGrpc.Mock.On("ReduceStocks", mock.Anything, data).Return(returnArg1)
}

func (c *CreateOrderTestSuite) CreateTransactionReq() *dto.TransactionReq {
	orderId, _ := gonanoid.New()

	return &dto.TransactionReq{
		Order: &entity.Order{
			OrderId:         orderId,
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
			UserId:          "hyfca5Sq7nQcaY6ACksXP",
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
		},
	}
}

func TestRepository_CreateOrder(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}
