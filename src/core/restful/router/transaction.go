package router

import (
	"github.com/dwprz/prasorganic-order-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Transaction(app *fiber.App, h *handler.Transaction, m *middleware.Middleware) {
		// all
		app.Add("POST", "/api/orders/transactions", m.VerifyJwt, h.Create)
}