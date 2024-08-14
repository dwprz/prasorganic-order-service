package helper

import (
	"encoding/base64"
	"fmt"

	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
)

func FormatMidtransTxReq(data *dto.TransactionReq) *dto.MidtransTxReq {
	return &dto.MidtransTxReq{
		TransactionDetails: struct {
			OrderId     string "json:\"order_id\""
			GrossAmount int    "json:\"gross_amount\""
		}{
			OrderId:     data.Order.OrderId,
			GrossAmount: data.Order.GrossAmount,
		},

		CreditCard: struct {
			Secure bool "json:\"secure\""
		}{
			Secure: true,
		},

		CustomerDetails: struct {
			CustomerName string "json:\"customer_name\""
			Whatsapp     string "json:\"whatsapp\""
		}{
			CustomerName: data.Order.Buyer,
			Whatsapp:     data.Order.WhatsApp,
		},

		Callbacks: struct {
			Finish  string "json:\"finish\""
			Error   string "json:\"error\""
			Pending string "json:\"pending\""
		}{
			Finish:  fmt.Sprintf("%s/order-status?orderId=%s", config.Conf.FrontEnd.BaseUrl, data.Order.OrderId),
			Error:   fmt.Sprintf("%s/order-status?orderId=%s", config.Conf.FrontEnd.BaseUrl, data.Order.OrderId),
			Pending: fmt.Sprintf("%s/order-status?orderId=%s", config.Conf.FrontEnd.BaseUrl, data.Order.OrderId),
		},
	}
}

func CreateMidtransBasicAuth() (auth string) {
	str := base64.StdEncoding.EncodeToString([]byte(config.Conf.Midtrans.ServerKey))
	auth = fmt.Sprintf("Basic %s:", str)
	return auth
}
