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

func fullCTE(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {
	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) AS total_orders FROM orders
	),
	cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		ORDER BY
			created_at DESC
		LIMIT 
			$1 OFFSET $2
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
	SELECT cto.total_orders , co.* FROM cte_total_orders AS cto CROSS JOIN cte_orders AS co;
	`

	res := postgres.WithContext(ctx).Raw(query, limit, offset).Scan(&queryRes)
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

func fullCTE_WithJsonAgg(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {
	queryRes := new(entity.QueryJsonWithCount)

	query := `
	WITH cte_total_orders AS (
		SELECT COUNT(*) FROM orders
	),
	cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		ORDER BY
			created_at DESC
		LIMIT 
			$1 OFFSET $2
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

	res := postgres.WithContext(ctx).Raw(query, limit, offset).Scan(&queryRes)
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

	orders, _ := helper.FormatOrderWithProducts(dummyOrders)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: queryRes.Total,
	}, nil
}

func nonFullCTE_1(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders struct {
		Count int
	}

	if err := postgres.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders;`).Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		ORDER BY
			created_at DESC
		LIMIT 
			$1 OFFSET $2
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

	res := postgres.WithContext(ctx).Raw(query, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders, _ := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: totalOrders.Count,
	}, nil
}

func nonFullCTE_2(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders struct {
		Count int
	}

	if err := postgres.WithContext(ctx).Raw(`SELECT COUNT(*) FROM orders;`).Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	var orderIds []string

	if err := postgres.WithContext(ctx).Raw(`SELECT order_id FROM orders ORDER BY created_at ASC LIMIT $1 OFFSET $2;`, limit, offset).Scan(&orderIds).Error; err != nil {
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

	orders, _ := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: totalOrders.Count,
	}, nil
}

func gorm_1(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders int64

	if err := postgres.WithContext(ctx).Table("orders").Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	var queryRes []*entity.QueryJoin

	query := `
	WITH cte_order_ids AS (
		SELECT 
			order_id 
		FROM 
			orders 
		ORDER BY
			created_at DESC
		LIMIT 
			$1 OFFSET $2
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

	res := postgres.WithContext(ctx).Raw(query, limit, offset).Scan(&queryRes)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(queryRes) == 0 {
		return nil, &errors.Response{HttpCode: 404, Message: "orders not found"}
	}

	orders, _ := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: int(totalOrders),
	}, nil
}

func gorm_2(ctx context.Context, limit, offset int) (*dto.OrdersWithCountRes, error) {

	var totalOrders int64

	if err := postgres.WithContext(ctx).Table("orders").Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	var orderIds []*struct {
		OrderId string
	}

	err := postgres.WithContext(ctx).Table("orders").Select("order_id").Order("created_at DESC").Limit(limit).Offset(offset).Scan(&orderIds).Error
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

	orders, _ := helper.FormatOrderWithProducts(queryRes)
	helper.OrderByCreatedAtDesc(orders)

	return &dto.OrdersWithCountRes{
		Orders:      orders,
		TotalOrders: int(totalOrders),
	}, nil
}

func Benchmark_CompareQueryCTE(b *testing.B) {
	b.Run("Full CTE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fullCTE(context.Background(), 20, 0)
		}
	})

	b.Run("Full CTE with json agg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fullCTE_WithJsonAgg(context.Background(), 20, 0)
		}
	})

	b.Run("Non Full CTE 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonFullCTE_1(context.Background(), 20, 0)
		}
	})

	b.Run("Non FUll CTE 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonFullCTE_2(context.Background(), 20, 0)
		}
	})

	b.Run("GORM 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gorm_1(context.Background(), 20, 0)
		}
	})

	b.Run("GORM 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gorm_2(context.Background(), 20, 0)
		}
	})
}

// 1 ms = 1.000.000 ns
// 1 s = 1000 ms

//================================ Full CTE ================================
// test 1:
// Benchmark_CompareQueryCTE/Full_CTE-12               1629            742773 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.274s

// test 2:
// Benchmark_CompareQueryCTE/Full_CTE-12               1464            755414 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.204s

// test 3:
// Benchmark_CompareQueryCTE/Full_CTE-12               1497            745242 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.212s

//================================ Full CTE With JSON Agg ================================
// test 1:
// Benchmark_CompareQueryCTE/Full_CTE_with_json_agg-12                 1126           1034334 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.287s

// test 2:
// Benchmark_CompareQueryCTE/Full_CTE_with_json_agg-12                  966           1167719 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.269s

// test 3:
// Benchmark_CompareQueryCTE/Full_CTE_with_json_agg-12                 1069           1160321 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.383s

//================================ Non FUll CTE 1 ================================

// test 1:
// Benchmark_CompareQueryCTE/Non_Full_CTE_1-12                 1407            810431 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.242s

// test 2:
// Benchmark_CompareQueryCTE/Non_Full_CTE_1-12                 1448            813744 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.277s

// test 3
// Benchmark_CompareQueryCTE/Non_Full_CTE_1-12                 1218            832710 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.129s

//================================ Non Full CTE 2 ================================

// test 1:
// Benchmark_CompareQueryCTE/Non_FUll_CTE_2-12                 1304            880243 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.256s

// test 2:
// Benchmark_CompareQueryCTE/Non_FUll_CTE_2-12                 1328            899630 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.301s

// test 3:
// Benchmark_CompareQueryCTE/Non_FUll_CTE_2-12                 1304            878167 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.253s

//================================ GORM 1 ================================

// test 1:
// Benchmark_CompareQueryCTE/GORM_1-12                 1400            816819 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.245s

// test 2:
// Benchmark_CompareQueryCTE/GORM_1-12                 1358            839860 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.244s

// test 3:
// Benchmark_CompareQueryCTE/GORM_1-12                 1392            829174 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.256s

//================================ GORM 2 ================================

// test 1:
// Benchmark_CompareQueryCTE/GORM_2-12                 1281            881933 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.239s

// test 2:
// Benchmark_CompareQueryCTE/GORM_2-12                 1264            876514 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     1.219s

// test 3:
// Benchmark_CompareQueryCTE/GORM_2-12                 1275            930213 ns/op
// PASS
// ok      github.com/dwprz/prasorganic-order-service/src/repository/benchmark     2.254s