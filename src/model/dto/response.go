package dto

import "github.com/dwprz/prasorganic-order-service/src/model/entity"

type MidtransTxRes struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

type TransactionRes struct {
	OrderId     string `json:"order_id"`
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

type OrdersWithCountRes struct {
	Orders      []*entity.OrderWithProducts `json:"orders"`
	TotalOrders int                         `json:"total_orders"`
}