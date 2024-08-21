package dto

import "github.com/dwprz/prasorganic-order-service/src/model/entity"

type MidtransTxReq struct {
	TransactionDetails struct {
		OrderId     string `json:"order_id"`
		GrossAmount int    `json:"gross_amount"`
	} `json:"transaction_details"`

	CreditCard struct {
		Secure bool `json:"secure"`
	} `json:"credit_card"`

	CustomerDetails struct {
		CustomerName string `json:"customer_name"`
		Whatsapp     string `json:"whatsapp"`
	} `json:"customer_details"`

	Callbacks struct {
		Finish  string `json:"finish"`
		Error   string `json:"error"`
		Pending string `json:"pending"`
	} `json:"callbacks"`
}

type TransactionReq struct {
	Order    *entity.Order          `json:"order" validate:"required"`
	Products []*entity.ProductOrder `json:"products" validate:"required,dive"`
}

type GetOrdersByCurrentUserReq struct {
	UserId string `json:"user_id" validate:"required,min=21,max=21"`
	Page   int    `json:"page" validate:"required,max=100"`
}

type GetOrdersReq struct {
	Status string `json:"status" validate:"omitempty,max=20"`
	Page   int    `json:"page" validate:"required,max=100"`
}

type CancelOrderReq struct {
	UserId  string `json:"user_id" validate:"required,min=21,max=21"`
	OrderId string `json:"order_id" validate:"required,min=21,max=21"`
}

type UpdateStatusReq struct {
	OrderId string `json:"order_id" validate:"required,min=21,max=21"`
	Status  string `json:"status" validation:"required,max=20"`
}

type AddShippingIdReq struct {
	OrderId    string `json:"order_id" validate:"required,min=21,max=21"`
	ShippingId string `json:"shipping_id" validation:"required"`
}
