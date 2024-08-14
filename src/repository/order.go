package repository

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"gorm.io/gorm"
)

type OrderImpl struct {
	db         *gorm.DB
	grpcClient *client.Grpc
}

func NewOrder(db *gorm.DB, gc *client.Grpc) repository.Order {
	return &OrderImpl{
		db:         db,
		grpcClient: gc,
	}
}

func (o *OrderImpl) Create(ctx context.Context, data *dto.TransactionReq) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("orders").Create(data.Order).Error; err != nil {
			return err
		}

		if err := tx.Table("product_orders").Create(data.Products).Error; err != nil {
			return err
		}

		if err := o.grpcClient.Product.UpdateStock(ctx, data.Products); err != nil {
			return err
		}

		return nil
	})

	return err
}
