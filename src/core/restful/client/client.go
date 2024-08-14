package client

import "github.com/dwprz/prasorganic-order-service/src/interface/delivery"

type Restful struct {
	Midtrans delivery.Midtrans
}

func NewRestful(md delivery.Midtrans) *Restful {
	return &Restful{
		Midtrans: md,
	}
}
