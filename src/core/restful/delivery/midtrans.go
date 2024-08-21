package delivery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/interface/delivery"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
)

type MidtransRESTfulImpl struct{}

func NewMidtransRESTful() delivery.MidtransRESTful {
	return &MidtransRESTfulImpl{}
}

func (m *MidtransRESTfulImpl) Transaction(ctx context.Context, data *dto.TransactionReq) (*dto.MidtransTxRes, error) {
	res, err := cbreaker.Midtrans.Execute(func() (any, error) {
		txReq := helper.FormatMidtransTxReq(data)

		uri := config.Conf.Midtrans.BaseUrl + "/snap/v1/transactions"

		a := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(a)

		auth := helper.CreateMidtransBasicAuth()
		b, err := json.Marshal(txReq)
		if err != nil {
			return nil, err
		}

		req := a.Request()
		req.Header.Set("Authorization", auth)
		req.Header.SetContentType("Application/json")
		req.Header.SetMethod("POST")
		req.SetBodyRaw([]byte(b))
		req.SetRequestURI(uri)

		if err := a.Parse(); err != nil {
			return nil, err
		}

		code, body, _ := a.Bytes()
		if code != 201 {
			return nil, &errors.Response{HttpCode: code, Message: string(body)}
		}

		res := new(dto.MidtransTxRes)
		helper.LogJSON(res)
		err = json.Unmarshal(body, res)

		return res, err
	})

	if err != nil {
		return nil, err
	}

	txRes, ok := res.(*dto.MidtransTxRes)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T expected *dto.MidtransTxRes", txRes)
	}

	return txRes, nil
}
