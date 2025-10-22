package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	MaxParticipants uint8 `json:"max_participants"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoomCreate struct {
	Name string `json:"name"`
	Status string `json:"status"`
	MaxParticipants uint8 `json:"max_participants"`
	UserID uuid.UUID `json:"user_id"`
}