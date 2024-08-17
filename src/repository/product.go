package repository

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"gorm.io/gorm"
)

type ProductImpl struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) repository.Product {
	return &ProductImpl{
		db: db,
	}
}

func (p *ProductImpl) FindByOrderId(ctx context.Context, orderId string) ([]*entity.ProductOrder, error) {
	var products []*entity.ProductOrder
	if err := p.db.WithContext(ctx).Table("product_orders").Where("order_id = ?", orderId).Scan(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "product order not found"}
	}

	return products, nil
}
