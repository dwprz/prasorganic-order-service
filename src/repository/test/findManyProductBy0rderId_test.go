package test

import (
	"context"
	"testing"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	repoimpl "github.com/dwprz/prasorganic-order-service/src/repository"
	"github.com/dwprz/prasorganic-order-service/test/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// go test -p=1 -v ./src/repository/test/... -count=1
// go test -run ^TestRepository_FindManyProductByOrderId$ -v ./src/repository/test/ -count=1

type FindManyProductByOrderIdTestSuite struct {
	suite.Suite
	order         *entity.OrderWithProducts
	productRepo   repository.Product
	postgresDB    *gorm.DB
	orderTestUtil *util.OrderTest
}

func (f *FindManyProductByOrderIdTestSuite) SetupSuite() {
	f.postgresDB = database.NewPostgres()
	f.productRepo = repoimpl.NewProduct(f.postgresDB)

	f.orderTestUtil = util.NewOrderTest(f.postgresDB)
	f.order = f.orderTestUtil.Create()
}

func (f *FindManyProductByOrderIdTestSuite) TearDownSuite() {
	f.orderTestUtil.Delete()

	sqlDB, err := f.postgresDB.DB()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.FindManyProductByOrderIdTestSuite/TearDownSuite", "section": "postgresDB.DB"}).Fatal(err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.FindManyProductByOrderIdTestSuite/TearDownSuite", "section": "sqlDB.Close"}).Fatal(err)
	}
}

func (f *FindManyProductByOrderIdTestSuite) Test_Success() {
	res, err := f.productRepo.FindByOrderId(context.Background(), f.order.Order.OrderId)
	assert.NoError(f.T(), err)

	assert.NotEmpty(f.T(), res)
}

func (f *FindManyProductByOrderIdTestSuite) Test_NotFound() {
	orderId := "not-found"

	res, err := f.productRepo.FindByOrderId(context.Background(), orderId)
	assert.Error(f.T(), err)

	resErr, ok := err.(*errors.Response)
	assert.True(f.T(), ok)

	assert.Equal(f.T(), 404, resErr.HttpCode)
	assert.Nil(f.T(), res)
}

func TestRepository_FindManyProductByOrderId(t *testing.T) {
	suite.Run(t, new(FindManyProductByOrderIdTestSuite))
}
