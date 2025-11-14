package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/talkey-be/handler"
	"github.com/heronhoga/talkey-be/util/auth"
)

func RegisterUserRoutes(app fiber.Router, userHandler *handler.UserHandler) {
	users := app.Group("/users")
	users.Post("/register", userHandler.RegisterNewUser)
	users.Post("/login", userHandler.Login)
	users.Post("/resetpassword", auth.AuthMiddleware, userHandler.ResetPassword)
}
