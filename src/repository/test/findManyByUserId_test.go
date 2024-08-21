package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/mock/delivery"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	repoimpl "github.com/dwprz/prasorganic-order-service/src/repository"
	"github.com/dwprz/prasorganic-order-service/test/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// go test -p=1 -v ./src/repository/test/... -count=1
// go test -run ^TestRepository_FindManyByUserId$ -v ./src/repository/test/ -count=1

type FindManyByUserIdTestSuite struct {
	suite.Suite
	order         *entity.OrderWithProducts
	orderRepo     repository.Order
	postgresDB    *gorm.DB
	orderTestUtil *util.OrderTest
}

func (f *FindManyByUserIdTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	productGrpc := delivery.NewProductGrpcMock()
	productConn := new(grpc.ClientConn)

	grpcClient := client.NewGrpc(productGrpc, productConn)

	f.orderRepo = repoimpl.NewOrder(f.postgresDB, grpcClient)

	f.orderTestUtil = util.NewOrderTest(f.postgresDB)

	f.order = f.orderTestUtil.Create()
}

func (f *FindManyByUserIdTestSuite) TearDownSuite() {
	f.orderTestUtil.Delete()

	sqlDB, err := f.postgresDB.DB()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.FindManyByUserIdTestSuite/TearDownSuite", "section": "postgresDB.DB"}).Fatal(err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.FindManyByUserIdTestSuite/TearDownSuite", "section": "sqlDB.Close"}).Fatal(err)
	}
}

func (f *FindManyByUserIdTestSuite) Test_Success() {
	limit := 20
	offset := 0

	res, err := f.orderRepo.FindManyByUserId(context.Background(), f.order.Order.UserId, limit, offset)
	assert.NoError(f.T(), err)

	assert.NotEmpty(f.T(), res)
}

func (f *FindManyByUserIdTestSuite) Test_NotFound() {
	userId := "not-found"
	limit := 20
	offset := 100

	res, err := f.orderRepo.FindManyByUserId(context.Background(), userId, limit, offset)
	assert.Error(f.T(), err)

	resErr, ok := err.(*errors.Response)
	assert.True(f.T(), ok)

	assert.Equal(f.T(), 404, resErr.HttpCode)
	assert.Nil(f.T(), res)
}

func TestRepository_FindManyByUserId(t *testing.T) {
	suite.Run(t, new(FindManyByUserIdTestSuite))
}
