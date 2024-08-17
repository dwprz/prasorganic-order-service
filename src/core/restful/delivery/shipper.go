package delivery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dwprz/prasorganic-order-service/src/common/helper"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-order-service/src/interface/delivery"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/gofiber/fiber/v2"
)

type ShipperImpl struct{}

func NewShipper() delivery.Shipper {
	return &ShipperImpl{}
}

func (s *ShipperImpl) ShippingOrder(ctx context.Context, data *entity.OrderWithProducts) (shippingId string, err error) {
	res, err := cbreaker.Shipper.Execute(func() (any, error) {
		shippingOrder := helper.FormatShippingOrderReq(data)

		uri := config.Conf.Shipper.BaseUrl + "/v3/order"

		a := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(a)

		a.JSON(shippingOrder)

		req := a.Request()
		req.Header.SetContentType("application/json")
		req.Header.Set("X-API-KEY", config.Conf.Shipper.ApiKey)
		req.Header.SetMethod("POST")
		req.SetRequestURI(uri)

		if err := a.Parse(); err != nil {
			return "", err
		}

		code, body, _ := a.Bytes()
		if code != 201 {
			return "", fmt.Errorf(string(body))
		}

		res := new(struct {
			Data struct {
				ShippingId string `json:"order_id"`
			} `json:"data"`
		})

		err = json.Unmarshal(body, res)

		return res.Data.ShippingId, err
	})

	return res.(string), err
}
