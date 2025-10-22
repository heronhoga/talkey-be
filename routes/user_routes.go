package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/talkey-be/handler"
)

func RegisterUserRoutes(app fiber.Router, userHandler *handler.UserHandler) {
	users := app.Group("/users")
	users.Post("/register", userHandler.RegisterNewUser)
	users.Post("/login", userHandler.Login)
}
