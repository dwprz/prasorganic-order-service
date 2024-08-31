package helper

import (
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

func FormatOrderWithProducts(data []*entity.QueryJoin) (orders []*entity.OrderWithProducts, totalOrders int) {

	dummy := make(map[string]*entity.OrderWithProducts)
	for _, order := range data {

		if dummy[order.OrderId] == nil {
			dummy[order.OrderId] = &entity.OrderWithProducts{
				Order: &entity.Order{
					OrderId:         order.OrderId,
					GrossAmount:     order.GrossAmount,
					Status:          order.Status,
					ShippingId:      order.ShippingId,
					Courier:         order.Courier,
					RateId:          order.RateId,
					RateName:        order.RateName,
					RateType:        order.RateType,
					COD:             order.COD,
					UseInsurance:    order.UseInsurance,
					PackageType:     order.PackageType,
					PaymentMethod:   order.PaymentMethod,
					SnapToken:       order.SnapToken,
					SnapRedirectURL: order.SnapRedirectURL,
					UserId:          order.UserId,
					Email:           order.Email,
					Buyer:           order.Buyer,
					Height:          order.Height,
					Length:          order.Length,
					Width:           order.Width,
					Weight:          order.Weight,
					AddressOwner:    order.AddressOwner,
					Street:          order.Street,
					AreaId:          order.AreaId,
					Area:            order.Area,
					Lat:             order.Lat,
					Lng:             order.Lng,
					Suburb:          order.Suburb,
					City:            order.City,
					Province:        order.Province,
					WhatsApp:        order.WhatsApp,
					CreatedAt:       order.CreatedAt,
					UpdatedAt:       order.UpdatedAt,
				},
			}
		}

		dummy[order.OrderId].Products = append(dummy[order.OrderId].Products, &entity.ProductOrder{
			OrderId:     order.OrderId,
			ProductId:   order.ProductId,
			ProductName: order.ProductName,
			Image:       order.Image,
			Price:       order.Price,
			Quantity:    order.Quantity,
		})
	}

	return MapValues(dummy), data[0].TotalOrders
}
