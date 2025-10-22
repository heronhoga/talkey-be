package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/heronhoga/talkey-be/model"
	"github.com/jackc/pgx/v5"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, roomCreate *model.RoomCreate, userId string) error
	JoinRoom(ctx context.Context, roomId string, userId string) error
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
		return errors.New("failed to begin the transaction")
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
		return errors.New("failed to insert the room")
	}

	createRoomParticipantQuery := `
		INSERT INTO room_participants (room_id, user_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, createRoomParticipantQuery, roomID, userId)
	if err != nil {
		return errors.New("failed to insert the room participant (create new room)")
	}

	return nil
}

func (r *roomRepository) JoinRoom(ctx context.Context, roomId string, userId string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.New("failed to begin the transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var existingRoomID string
	checkRoomQuery := `SELECT id FROM rooms WHERE status = '1' AND id = $1 LIMIT 1`
	err = tx.QueryRow(ctx, checkRoomQuery, roomId).Scan(&existingRoomID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return errors.New("the room is not available")
		}
		return errors.New("failed to check the room")
	}

	var existingUserID string
	checkUserExistsQuery := `SELECT id FROM room_participants WHERE room_id = $1 AND user_id = $2`
	err = tx.QueryRow(ctx, checkUserExistsQuery, roomId, userId).Scan(&existingUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
		} else {
			// actual database error
			return errors.New("failed to check existing user")
		}
	} else {
		return errors.New("user already exists in the room")
	}

	addParticipantQuery := `INSERT INTO room_participants (room_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, addParticipantQuery, roomId, userId)
	if err != nil {
		return errors.New("failed to add new participant")
	}

	return nil
}

