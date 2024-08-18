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
	Order    *entity.Order          `json:"order"`
	Products []*entity.ProductOrder `json:"products"`
}

type GetOrdersByCurrentUserReq struct {
	UserId string `validate:"required,min=21,max=21"`
	Page   int    `validate:"required,max=100"`
}