package handler

import (
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Transaction struct {
	txService service.Transaction
}

func NewTransaction(ts service.Transaction) *Transaction {
	return &Transaction{
		txService: ts,
	}
}

func (t *Transaction) Create(c *fiber.Ctx) error {
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
