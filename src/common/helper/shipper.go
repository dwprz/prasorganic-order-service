package helper

import (
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
)

func FormatShippingOrderReq(data *entity.OrderWithProducts) *entity.ShippingOrder {
	var items []entity.Item
	for _, product := range data.Products {

		items = append(items, entity.Item{
			Name:  product.ProductName,
			Price: product.Price,
			Qty:   product.Quantity,
		})
	}

	return &entity.ShippingOrder{
		Consignee: entity.Consignee{
			Name:        data.Order.AddressOwner,
			PhoneNumber: data.Order.WhatsApp,
		},
		Consigner: entity.Consigner{
			Name:        config.Conf.Store.Name,
			PhoneNumber: config.Conf.Store.PhoneNumber,
		},
		Courier: entity.Courier{
			COD:          data.Order.COD,
			RateId:       data.Order.RateId,
			UseInsurance: data.Order.UseInsurance,
		},
		Coverage: config.Conf.Shipping.Coverage,
		Destination: entity.Destination{
			Address: data.Order.Street,
			AreaId:  data.Order.AreaId,
			Lat:     data.Order.Lat,
			Lng:     data.Order.Lng,
		},
		ExternalId: data.Order.OrderId,
		Origin: entity.Origin{
			Address: config.Conf.Store.Address,
			AreaId:  config.Conf.Store.AreaId,
			Lat:     config.Conf.Store.Latitude,
			Lng:     config.Conf.Store.Longitude,
		},
		Package: entity.Package{
			Height:      data.Order.Height,
			Length:      data.Order.Length,
			PackageType: data.Order.PackageType,
			Price:       data.Order.GrossAmount,
			Weight:      data.Order.Weight,
			Width:       data.Order.Width,
			Items:       items,
		},
		PaymentType: config.Conf.Shipping.PaymentType,
	}
}
