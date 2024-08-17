package client

import "github.com/dwprz/prasorganic-order-service/src/interface/delivery"

type Restful struct {
	Midtrans delivery.Midtrans
	Shipper  delivery.Shipper
}

func NewRestful(md delivery.Midtrans, sd delivery.Shipper) *Restful {
	return &Restful{
		Midtrans: md,
		Shipper:  sd,
	}
}
