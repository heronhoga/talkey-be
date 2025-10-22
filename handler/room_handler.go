package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/talkey-be/model"
	"github.com/heronhoga/talkey-be/service"
)

type RoomHandler struct {
	service *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: roomService,
	}
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req *model.RoomCreate

	userId := c.Locals("id").(string)


	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.service.CreateRoom(c.Context(), req, userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}