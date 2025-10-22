package service

import (
	"context"
	"errors"

	"github.com/heronhoga/talkey-be/model"
	"github.com/heronhoga/talkey-be/repository"
)

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(ctx context.Context, roomCreate *model.RoomCreate, userId string) error {
	if roomCreate.Name == "" || roomCreate.Status == "" || roomCreate.MaxParticipants == 0 {
		return errors.New("name, status, and max participants are required")
	}

	if roomCreate.Status != "1" && roomCreate.Status != "0" {
		return errors.New("room status must be between 0 and 1")
	}

	if roomCreate.MaxParticipants > 8 || roomCreate.MaxParticipants < 0 {
		return errors.New("the number of participants exceed the limit")
	}

	//room model
	newRoom := &model.RoomCreate{
		Name: roomCreate.Name,
		Status: roomCreate.Status,
		MaxParticipants: roomCreate.MaxParticipants,
	}

	return s.repo.CreateRoom(ctx, newRoom, userId)
}