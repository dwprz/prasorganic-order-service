package handler

import (
	"strconv"

	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Order struct {
	txService    service.Transaction
	orderService service.Order
}

func NewOrder(ts service.Transaction, os service.Order) *Order {
	return &Order{
		txService:    ts,
		orderService: os,
	}
}

func (t *Order) Transaction(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	userId := userData["user_id"].(string)

	req := new(dto.TransactionReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Order.UserId = userId
	res, err := t.txService.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{"data": res})
}

func (t *Order) GetByCurrentUser(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	userId := userData["user_id"].(string)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return err
	}

	res, err := t.orderService.FindManyByUserId(c.Context(), &dto.GetOrdersByCurrentUserReq{
		UserId: userId,
		Page:   page,
	})

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": res.Data, "paging": res.Paging})
}

func (t *Order) Get(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return err
	}

	status := c.Query("status", "paid")

	res, err := t.orderService.FindMany(c.Context(), &dto.GetOrdersReq{
		Status: status,
		Page:   page,
	})

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": res.Data, "paging": res.Paging})
}

func (t *Order) Cancellation(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)
	userId := userData["user_id"].(string)

	orderId := c.Params("orderId")

	err := t.orderService.Cancel(c.Context(), &dto.CancelOrderReq{
		UserId:  userId,
		OrderId: orderId,
	})

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": "cancelled order successfully"})
}

func (t *Order) UpdateStatus(c *fiber.Ctx) error {
	req := new(dto.UpdateStatusReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.OrderId = c.Params("orderId")

	err := t.orderService.UpdateStatus(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"data": "updated order status successfully"})
}