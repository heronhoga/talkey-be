package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/talkey-be/handler"
	"github.com/heronhoga/talkey-be/util/auth"
)

func RegisterRoomRoutes(app fiber.Router, roomHandler *handler.RoomHandler) {
	users := app.Group("/rooms")
	users.Post("/create", auth.AuthMiddleware, roomHandler.CreateRoom)
	users.Post("/join", auth.AuthMiddleware, roomHandler.JoinRoom)
}