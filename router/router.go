package router

import (
	"github.com/NishantBansal2003/Brokerax/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// app.Post("/api/auth/login", controller.Login)
	// app.Post("/api/auth/signup", controller.Signup)

	// app.Post("/api/user/portfolio", controller.Portfolio)
	// app.Post("/api/user/stock/add", controller.AddStock)
	// app.Post("/api/user/stock/remove", controller.RemoveStock)
	
	// Grouping authentication routes
	auth := app.Group("/api/auth")
	auth.Post("/login", controller.Login)
	auth.Post("/signup", controller.Signup)

	// Grouping user routes
	user := app.Group("/api/user")
	user.Post("/portfolio", controller.Portfolio)
	user.Post("/stock/add", controller.AddStock)
	user.Post("/stock/remove", controller.RemoveStock)
}
