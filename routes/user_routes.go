package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/talkey-be/handler"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(app fiber.Router, userHandler *handler.UserHandler) {
	users := app.Group("/users")

	users.Get("/register", userHandler.RegisterNewUser)
}
