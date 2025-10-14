package service

import (
	"context"
	"errors"
	"regexp"

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

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email format is not valid")
	}

	// Validate password length
	if len(password) < 6 {
		return errors.New("minimum password length is 6 characters")
	}

	// Validate existing user email
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
