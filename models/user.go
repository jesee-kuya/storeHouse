package models

import (
	"time"
)

// User represents a system user with different roles and permissions
type User struct {
	ID          string    `json:"id" db:"id"`
	Username    string    `json:"username" db:"username" binding:"required,max=50"`
	Email       string    `json:"email" db:"email" binding:"required,email"`
	PasswordHash string   `json:"-" db:"password_hash" binding:"required"`
	FullName    string    `json:"full_name" db:"full_name" binding:"required,max=200"`
	Role        string    `json:"role" db:"role" binding:"required"`
	PhoneNumber string    `json:"phone_number" db:"phone_number" binding:"max=12"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastLogin   *time.Time `json:"last_login" db:"last_login"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole represents the different user roles
type UserRole string

const (
	RoleAdmin     UserRole = "Admin"
	RoleTreasurer UserRole = "Treasurer"
	RoleClerk     UserRole = "Clerk"
)

// ValidateRole checks if the user role is valid
func (u *User) ValidateRole() error {
	switch u.Role {
	case string(RoleAdmin), string(RoleTreasurer), string(RoleClerk):
		return nil
	default:
		return ErrInvalidRole
	}
}

// CreateUserRequest represents the request for creating a new user
type CreateUserRequest struct {
	Username    string    `json:"username" binding:"required,max=50"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=8"`
	FullName    string    `json:"full_name" binding:"required,max=200"`
	Role        string    `json:"role" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"max=12"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Username    *string   `json:"username" binding:"max=50"`
	Email       *string   `json:"email" binding:"email"`
	FullName    *string   `json:"full_name" binding:"max=200"`
	Role        *string   `json:"role"`
	PhoneNumber *string   `json:"phone_number" binding:"max=12"`
	IsActive    *bool     `json:"is_active"`
}

// UserResponse represents the user response without sensitive data
type UserResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	LastLogin   *time.Time `json:"last_login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		FullName:    u.FullName,
		Role:        u.Role,
		PhoneNumber: u.PhoneNumber,
		IsActive:    u.IsActive,
		LastLogin:   u.LastLogin,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}