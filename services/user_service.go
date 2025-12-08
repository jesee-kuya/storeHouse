package services

import (
	"errors"
	"regexp"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	DB *sqlx.DB
}

// Create a new instance of UserService
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser handles user creation business logic
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.UserResponse, error) {
	// Validate password strength
	if err := s.validatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check for duplicate username
	if _, err := repository.GetUserByUsername(s.DB, req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check for duplicate email
	if _, err := repository.GetUserByEmail(s.DB, req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Store password (in production, this should be hashed)
	passwordHash := req.Password

	// Validate user role
	user := models.User{Role: req.Role}
	if err := user.ValidateRole(); err != nil {
		return nil, err
	}

	// Prepare model for DB
	userModel := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FullName:     req.FullName,
		Role:         req.Role,
		PhoneNumber:  req.PhoneNumber,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to DB
	newUser, err := repository.CreateUser(s.DB, userModel)
	if err != nil {
		return nil, err
	}

	return newUser.ToResponse(), nil
}

// UpdateUser handles update logic
func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	// Fetch existing record
	existing, err := repository.GetUser(s.DB, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Apply updates only if fields are provided
	if req.Username != nil {
		// Check for duplicate username if different
		if existing.Username != *req.Username {
			if _, err := repository.GetUserByUsername(s.DB, *req.Username); err == nil {
				return nil, errors.New("username already exists")
			}
		}
		existing.Username = *req.Username
	}
	if req.Email != nil {
		// Check for duplicate email if different
		if existing.Email != *req.Email {
			if _, err := repository.GetUserByEmail(s.DB, *req.Email); err == nil {
				return nil, errors.New("email already exists")
			}
		}
		existing.Email = *req.Email
	}
	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.Role != nil {
		existing.Role = *req.Role
		// Validate user role
		if err := existing.ValidateRole(); err != nil {
			return nil, err
		}
	}
	if req.PhoneNumber != nil {
		existing.PhoneNumber = *req.PhoneNumber
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateUser(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeactivateUser sets user as inactive
func (s *UserService) DeactivateUser(id string) error {
	// Ensure exists before deactivating
	if _, err := repository.GetUser(s.DB, id); err != nil {
		return errors.New("user not found")
	}

	_, err := repository.DeactivateUser(s.DB, id)
	return err
}

// DeleteUser removes a user record
func (s *UserService) DeleteUser(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetUser(s.DB, id); err != nil {
		return errors.New("user not found")
	}

	_, err := s.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// GetUser returns single user details
func (s *UserService) GetUser(id string) (*models.UserResponse, error) {
	user, err := repository.GetUser(s.DB, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

// GetUserByUsername returns user by username
func (s *UserService) GetUserByUsername(username string) (*models.UserResponse, error) {
	user, err := repository.GetUserByUsername(s.DB, username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

// GetUserByEmail returns user by email
func (s *UserService) GetUserByEmail(email string) (*models.UserResponse, error) {
	user, err := repository.GetUserByEmail(s.DB, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := repository.GetAllUsers(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *u.ToResponse())
	}

	return responses, nil
}

// GetUsersByRole returns users with a specific role
func (s *UserService) GetUsersByRole(role string) ([]models.UserResponse, error) {
	// Validate user role
	user := models.User{Role: role}
	if err := user.ValidateRole(); err != nil {
		return nil, err
	}

	users, err := repository.GetUsersByRole(s.DB, role)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *u.ToResponse())
	}

	return responses, nil
}

// GetActiveUsers returns all active users
func (s *UserService) GetActiveUsers() ([]models.UserResponse, error) {
	users, err := repository.GetActiveUsers(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *u.ToResponse())
	}

	return responses, nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (s *UserService) UpdateLastLogin(id string) error {
	// Ensure user exists
	if _, err := repository.GetUser(s.DB, id); err != nil {
		return errors.New("user not found")
	}

	return repository.UpdateLastLogin(s.DB, id)
}

// AuthenticateUser validates user credentials
func (s *UserService) AuthenticateUser(username, password string) (*models.UserResponse, error) {
	// Get user by username
	user, err := repository.GetUserByUsername(s.DB, username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password (simple comparison - use proper hashing in production)
	if user.PasswordHash != password {
		return nil, errors.New("invalid username or password")
	}

	// Update last login
	if err := repository.UpdateLastLogin(s.DB, user.ID); err != nil {
		// Don't return error for last login update failure
		// Just log it or handle it as a warning
	}

	return user.ToResponse(), nil
}

// ChangePassword changes user's password
func (s *UserService) ChangePassword(id string, oldPassword, newPassword string) error {
	// Get user
	user, err := repository.GetUser(s.DB, id)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password (simple comparison - use proper hashing in production)
	if user.PasswordHash != oldPassword {
		return errors.New("invalid old password")
	}

	// Validate new password strength
	if err := s.validatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Update password (simple storage - use proper hashing in production)
	_, err = s.DB.Exec("UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3", newPassword, time.Now(), id)
	return err
}

// validatePasswordStrength validates password meets security requirements
func (s *UserService) validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	return nil
}
