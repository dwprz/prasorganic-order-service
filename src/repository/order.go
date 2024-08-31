package repository

import (
	"context"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	"github.com/dwprz/prasorganic-order-service/src/core/grpc/client"
	"github.com/dwprz/prasorganic-order-service/src/interface/repository"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"google.golang.org/grpc/codes"
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

		if err := o.grpcClient.Product.ReduceStocks(ctx, data.Products); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (o *OrderImpl) FindById(ctx context.Context, orderId string) (*entity.OrderWithProducts, error) {
	var queryRes []*entity.QueryJoin

	query := `
	SELECT 
		* 
	FROM 
		orders AS o 
	INNER JOIN 
		product_orders AS po ON o.order_id = po.order_id
	WHERE
		o.order_id = $1;	
	`

	res := o.db.WithContext(ctx).Raw(query, orderId).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "order not found"}
	}

	orders, _ := helper.FormatOrderWithProducts(queryRes)

	return orders[0], nil
}

func (o *OrderImpl) FindMany(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {
	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) AS total_orders FROM orders
	),
	cte_order_ids AS (
		SELECT order_id FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2
	), 
	cte_orders AS (
		SELECT 
			*
		FROM 
			orders AS o 
		INNER JOIN 
			product_orders AS po ON o.order_id = po.order_id
		WHERE
			o.order_id IN (SELECT order_id FROM cte_order_ids)
	)
	SELECT cto.total_orders, co.* FROM cte_total_orders AS cto CROSS JOIN cte_orders AS co;
	`

	res := o.db.WithContext(ctx).Raw(query, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders, total := helper.FormatOrderWithProducts(queryRes)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: total,
	}, nil
}

func (o *OrderImpl) FindManyByUserId(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {
	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) AS total_orders FROM orders WHERE user_id = $1
	),
	cte_order_ids AS (
		SELECT order_id FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	), 
	cte_orders AS (
		SELECT 
			*
		FROM 
			orders AS o 
		INNER JOIN 
			product_orders AS po ON o.order_id = po.order_id
		WHERE
			o.order_id IN (SELECT order_id FROM cte_order_ids)
	)
	SELECT cto.total_orders, co.* FROM cte_total_orders AS cto CROSS JOIN cte_orders AS co;
	`

	res := o.db.WithContext(ctx).Raw(query, userId, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders, total := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: total,
	}, nil
}

func (o *OrderImpl) FindManyByStatus(ctx context.Context, status string, limit, offset int) (*dto.OrdersWithCountRes, error) {
	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) AS total_orders FROM orders WHERE status = $1
	),
	cte_order_ids AS (
		SELECT order_id FROM orders WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	), 
	cte_orders AS (
		SELECT 
			*
		FROM 
			orders AS o 
		INNER JOIN 
			product_orders AS po ON o.order_id = po.order_id
		WHERE
			o.order_id IN (SELECT order_id FROM cte_order_ids)
	)
	SELECT cto.total_orders, co.* FROM cte_total_orders AS cto CROSS JOIN cte_orders AS co;
	`

	res := o.db.WithContext(ctx).Raw(query, status, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders, total := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: total,
	}, nil
}

func (o *OrderImpl) UpdateById(ctx context.Context, data *entity.Order) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		res := o.db.Table("orders").Updates(data)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return &errors.Response{HttpCode: 404, GrpcCode: codes.NotFound, Message: "order not found"}
		}

		return nil
	})

	return err
}
