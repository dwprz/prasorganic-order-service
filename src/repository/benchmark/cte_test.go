package benchmark

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/database"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"gorm.io/gorm"
)

// *cd current directory
// go test -v -bench=. -count=1 -p=1

var postgres *gorm.DB

func init() {
	postgres = database.NewPostgres()
}

func fullCTE(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {
	queryRes := new(entity.QueryJsonWithCount)

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) FROM orders WHERE user_id = $1
	),
	cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		WHERE 
			user_id = $1 
		ORDER BY
			created_at DESC
		LIMIT 
			$2 OFFSET $3
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
	SELECT 
		(SELECT * FROM cte_total_orders) AS total,
		(SELECT json_agg(row_to_json(cte_orders.*)) FROM cte_orders) AS data;
	`

	res := postgres.WithContext(ctx).Raw(query, userId, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes.Data) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	var dummyOrders []*entity.QueryJoin
	if err := json.Unmarshal(queryRes.Data, &dummyOrders); err != nil {
		return nil, err
	}

	orders := helper.FormatOrderWithProducts(dummyOrders)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: queryRes.Total,
	}, nil
}

func nonFullCTE_1(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders struct {
		Count int
	}

	if err := postgres.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders WHERE user_id = $1`, userId).Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		WHERE 
			user_id = $1 
		ORDER BY
			created_at DESC
		LIMIT 
			$2 OFFSET $3
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
	SELECT * FROM cte_orders;
	`

	res := postgres.WithContext(ctx).Raw(query, userId, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: totalOrders.Count,
	}, nil
}

func nonFullCTE_2(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders struct {
		Count int
	}

	if err := postgres.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders WHERE user_id = $1`, userId).Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	var orderIds []string

	if err := postgres.WithContext(ctx).Raw(`SELECT order_id FROM orders WHERE user_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`, userId, limit, offset).Scan(&orderIds).Error; err != nil {
		return nil, err
	}

	var queryRes []*entity.QueryJoin

	query := `
	SELECT 
			*
	FROM 
		orders AS o 
	INNER JOIN 
		product_orders AS po ON o.order_id = po.order_id
	WHERE
		o.order_id IN (?);
	`

	res := postgres.WithContext(ctx).Raw(query, orderIds).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: totalOrders.Count,
	}, nil
}

func gorm_1(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders int64

	if err := postgres.WithContext(ctx).Table("orders").Where("user_id = ?", userId).Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		WHERE 
			user_id = $1 
		ORDER BY
			created_at DESC
		LIMIT 
			$2 OFFSET $3
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
	SELECT * FROM cte_orders;
	`

	res := postgres.WithContext(ctx).Raw(query, userId, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: int(totalOrders),
	}, nil
}

func gorm_2(ctx context.Context, userId string, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders int64

	if err := postgres.WithContext(ctx).Table("orders").Where("user_id = ?", userId).Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	var orderIds []*struct {
		OrderId string
	}

	err := postgres.WithContext(ctx).Table("orders").Select("order_id").Where("user_id = ?", userId).Order("created_at DESC").Limit(limit).Offset(offset).Scan(&orderIds).Error
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, v := range orderIds {
		ids = append(ids, v.OrderId)
	}

	var queryRes []*entity.QueryJoin

	query := `
	SELECT 
			*
	FROM 
		orders AS o 
	INNER JOIN 
		product_orders AS po ON o.order_id = po.order_id
	WHERE
		o.order_id IN (?);
	`

	res := postgres.WithContext(ctx).Raw(query, ids).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: int(totalOrders),
	}, nil
}

func Benchmark_CompareQueryCTE(b *testing.B) {
	b.Run("Full CTE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fullCTE(context.Background(), "user_1", 20, 0)
		}
	})

	b.Run("Non Full CTE 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonFullCTE_1(context.Background(), "user_1", 20, 0)
		}
	})

	b.Run("Non FUll CTE 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonFullCTE_2(context.Background(), "user_1", 20, 0)
		}
	})

	b.Run("GORM 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gorm_1(context.Background(), "user_1", 20, 0)
		}
	})

	b.Run("GORM 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gorm_2(context.Background(), "user_1", 20, 0)
		}
	})
}

// 1 ms = 1.000.000 ns
// 1 s = 1000 ms
//================================ Full CTE ================================
// test 1:
// Benchmark_CompareQueryCTE/With_CTE-12               2779            400939 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.173s

// test 2:
// Benchmark_CompareQueryCTE/With_CTE-12               2896            403285 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.052s

// test 3:
// Benchmark_CompareQueryCTE/With_CTE-12               2967            398861 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.227s

//================================ Non FUll CTE 1 ================================
// test 1:
// Benchmark_CompareQueryCTE/Non_CTE_1-12              2514            436977 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.162s

// test 2:
// Benchmark_CompareQueryCTE/Non_CTE_1-12              2533            438572 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.174s

// test 3:
// Benchmark_CompareQueryCTE/Non_CTE_1-12              2738            433828 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.223s

//================================ Non Full CTE 2 ================================
// test 1:
// Benchmark_CompareQueryCTE/Non_CTE_2-12              2257            516415 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.135s

// test 2:
// Benchmark_CompareQueryCTE/Non_CTE_2-12              2325            522300 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.080s

// test 3:
// Benchmark_CompareQueryCTE/Non_CTE_2-12              2163            516490 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.188s

//================================ GORM 1 ================================
// test 1:
// Benchmark_CompareQueryCTE/GORM_1-12                 2541            447598 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.200s

// test 2:
// Benchmark_CompareQueryCTE/GORM_1-12                 2512            450063 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.194s

// test 3:
// Benchmark_CompareQueryCTE/GORM_1-12                 2418            441662 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.135s

//================================ GORM 2 ================================
// test 1:
// Benchmark_CompareQueryCTE/GORM_2-12                 2143            530503 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.208s

// test 2:
// Benchmark_CompareQueryCTE/GORM_2-12                 2245            529915 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.174s

// test 3:
// Benchmark_CompareQueryCTE/GORM_2-12                 2073            541375 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.195s
