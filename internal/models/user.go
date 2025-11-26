package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // "-" means never return in JSON
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest - Request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
}

// LoginRequest - Request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse - Response after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

// UpdateUserRequest - Request body for updating user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=2"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}
