package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/heronhoga/talkey-be/model"
	"github.com/jackc/pgx/v5"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, roomCreate *model.RoomCreate, userId string) error
}

type roomRepository struct {
	db *pgx.Conn
}

func NewRoomRepository(db *pgx.Conn) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) CreateRoom(ctx context.Context, roomCreate *model.RoomCreate, userId string) error {

	//begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var roomID uuid.UUID
	createRoomQuery := `
		INSERT INTO rooms (name, status, max_participants)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(ctx, createRoomQuery, roomCreate.Name, roomCreate.Status, roomCreate.MaxParticipants).Scan(&roomID)
	if err != nil {
		return fmt.Errorf("failed to insert room: %w", err)
	}

	createRoomParticipantQuery := `
		INSERT INTO room_participants (room_id, user_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, createRoomParticipantQuery, roomID, userId)
	if err != nil {
		return fmt.Errorf("failed to insert room participant: %w", err)
	}

	return nil
}
