package helper

import (
	"math"

	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

func CreateLimitAndOffset(page int) (limit, offset int) {
	limit = 20
	offset = (page - 1) * limit

	return limit, offset
}

func FormatPagedData[T any](data T, totalData int, page int, limit int) *entity.DataWithPaging[T] {

	return &entity.DataWithPaging[T]{
		Data: data,
		Paging: &entity.Paging{
			TotalData: totalData,
			Page:      page,
			TotalPage: int(math.Ceil(float64(totalData) / float64(limit))),
		},
	}
}
