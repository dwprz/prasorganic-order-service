package router

import (
	"github.com/dwprz/prasorganic-order-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-order-service/src/core/restful/middleware"
	"github.com/gofiber/fiber/v2"
)

func Order(app *fiber.App, h *handler.OrderRESTful, m *middleware.Middleware) {
	// super admin
	app.Add("PATCH", "/api/orders/:orderId/statuses", m.VerifyJwt, m.VerifySuperAdmin, h.UpdateStatus)

	// admin & super admin
	app.Add("GET", "/api/orders", m.VerifyJwt, m.VerifyAdmin, h.Get)

	// all
	app.Add("POST", "/api/orders/transactions", m.VerifyJwt, h.Transaction)
	app.Add("GET", "/api/orders/users/current", m.VerifyJwt, h.GetByCurrentUser)
	app.Add("PATCH", "/api/orders/:orderId/cancellations", m.VerifyJwt, h.Cancellation)
}
