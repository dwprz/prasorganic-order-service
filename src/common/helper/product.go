package helper

import (
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	pb "github.com/dwprz/prasorganic-proto/protogen/product"
)

func ConvertToProductData(data []*entity.ProductOrder) []*pb.ProductData {
	var req []*pb.ProductData

	for _, v := range data {
		product := &pb.ProductData{
			ProductId: uint32(v.ProductId),
			Quantity:  uint32(v.Quantity),
		}
		req = append(req, product)
	}
	
	return req
}
