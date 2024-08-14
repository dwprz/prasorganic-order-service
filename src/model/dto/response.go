package dto

type MidtransTxRes struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

type TransactionRes struct {
	OrderId     string `json:"order_id"`
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}
