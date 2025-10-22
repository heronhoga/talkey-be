package model

import (
	"time"

	"github.com/google/uuid"
)

type RoomParticipants struct {
	ID uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
	LeftAt time.Time `json:"left_at"` 
}