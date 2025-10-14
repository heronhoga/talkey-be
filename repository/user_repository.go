package repository

import (
	"context"

	"github.com/heronhoga/talkey-be/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	CheckUserExists(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, username, email FROM users WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (r *userRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email=$1`
	var count int
	err := r.db.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
