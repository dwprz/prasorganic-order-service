package helper

import (
	"sort"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

func OrderByCreatedAtDesc(data []*entity.OrderWithProducts) {
	if len(data) > 1 {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Order.CreatedAt.After(data[j].Order.CreatedAt)
		})
	}
}
