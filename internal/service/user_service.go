package service

import (
	"apiserver/internal/models"
	"apiserver/internal/repository"
	"apiserver/internal/utils"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo:     repo,
		validate: validator.New(),
	}
}

// Register creates a new user account
func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns user data
func (s *UserService) Login(req *models.LoginRequest) (*models.User, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Get user from database
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Check if email is already taken by another user
		existingUser, _ := s.repo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}

		user.Email = req.Email
	}

	// Update in database
	if err := s.repo.Update(user); err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}

	// Clear password
	user.Password = ""
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
