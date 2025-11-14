package repository

import (
	"context"
	"errors"

	"github.com/heronhoga/talkey-be/model"
	"github.com/heronhoga/talkey-be/util/auth"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	CheckUserExists(ctx context.Context, email string) (bool, error)
	Login(ctx context.Context, user *model.UserLogin) (string, error)
	ResetPassword(ctx context.Context, resetPasswordRequest *model.UserResetPassword) error
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

func (r *userRepository) Login(ctx context.Context, user *model.UserLogin) (string, error) {
	query := `SELECT id, username, email, password 
	          FROM users 
	          WHERE username=$1 OR email=$1 
	          LIMIT 1`

	var existingUser model.User
	err := r.db.QueryRow(ctx, query, user.Username).Scan(
		&existingUser.ID,
		&existingUser.Username,
		&existingUser.Email,
		&existingUser.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No user found
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return "", errors.New("invalid username/email or password")
	}

	token, err := auth.GenerateToken(existingUser.ID.String(), existingUser.Username)
	return token, nil
}

func (r *userRepository) ResetPassword(ctx context.Context, resetPasswordRequest *model.UserResetPassword) error {
    // find existing user
    query := `SELECT password FROM users WHERE id = $1`
    var existingUserPassword string

    err := r.db.QueryRow(ctx, query, resetPasswordRequest.UserID).Scan(&existingUserPassword)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return errors.New("user not found")
        }
        return errors.New("failed to get user")
    }

    // compare password
    if err := bcrypt.CompareHashAndPassword(
        []byte(existingUserPassword),
        []byte(resetPasswordRequest.OldPassword),
    ); err != nil {
        return errors.New("old password is incorrect")
    }

    // hash new password
    hashedNewPassword, err := bcrypt.GenerateFromPassword(
        []byte(resetPasswordRequest.NewPassword),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return errors.New("failed to hash new password")
    }

    // update password
    updateQuery := `UPDATE users SET password = $1 WHERE id = $2`
    _, err = r.db.Exec(ctx, updateQuery, hashedNewPassword, resetPasswordRequest.UserID)
    if err != nil {
        return errors.New("failed to update password")
    }

    return nil
}


