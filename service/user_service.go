package service

import (
	"context"
	"errors"

	"github.com/heronhoga/talkey-be/model"
	"github.com/heronhoga/talkey-be/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) RegisterNewUser(ctx context.Context, username, email, password string) error {
	//validation
	if username == "" || email == "" || password == "" {
		return errors.New("username, email, and password are required")
	}

	//Check if user already exists
	exists, err := s.repo.CheckUserExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	//User model
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	//Save to database
	return s.repo.Create(ctx, user)
}
