package router

import (
	"github.com/NishantBansal2003/Brokerax/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/auth/login", controller.Login)
	app.Post("/api/auth/signup", controller.Signup)

	app.Post("/api/user/portfolio", controller.Portfolio)
	app.Post("/api/user//stock/add", controller.AddStock)
	app.Post("/api/user/stock/remove", controller.RemoveStock)
}
