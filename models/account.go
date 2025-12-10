package models

import (
	"errors"
	"time"
)

// Account represents transaction accounts
type Account struct {
	ID          string    `json:"id" db:"id"`
	AccountName string    `json:"account_name" db:"account_name" binding:"required,max=100"`
	AccountType string    `json:"account_type" db:"account_type" binding:"required"`
	LocalShare  *float64  `json:"local_share" db:"local_share"`
	Notes       *string   `json:"notes" db:"notes"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AccountType represents the different account types
type AccountType string

const (
	AccountBank      AccountType = "Bank"
	AccountExpense   AccountType = "Expense"
	AccountIncome    AccountType = "Income"
	AccountAsset     AccountType = "Asset"
	AccountLiability AccountType = "liability"
)

// ValidateAccountType checks if the account type is valid
func (a *Account) ValidateAccountType() error {
	switch a.AccountType {
	case string(AccountBank), string(AccountExpense), string(AccountIncome), string(AccountAsset), string(AccountLiability):
		return nil
	default:
		return ErrInvalidAccountType
	}
}

// CreateAccountRequest represents the request for creating a new account
type CreateAccountRequest struct {
	AccountName string   `json:"account_name" binding:"required,max=100"`
	AccountType string   `json:"account_type" binding:"required"`
	LocalShare  *float64 `json:"local_share"`
	Notes       *string  `json:"notes"`
}

// Validate validates the CreateAccountRequest
func (req *CreateAccountRequest) Validate() error {
	if req.AccountName == "" {
		return errors.New("account name is required")
	}
	if req.AccountType == "" {
		return errors.New("account type is required")
	}

	// Validate account type
	account := Account{AccountType: req.AccountType}
	return account.ValidateAccountType()
}

// UpdateAccountRequest represents the request for updating an account
type UpdateAccountRequest struct {
	AccountName *string  `json:"account_name" binding:"max=100"`
	AccountType *string  `json:"account_type"`
	LocalShare  *float64 `json:"local_share"`
	Notes       *string  `json:"notes"`
	IsActive    *bool    `json:"is_active"`
}

// AccountResponse represents the account response
type AccountResponse struct {
	ID          string    `json:"id"`
	AccountName string    `json:"account_name"`
	AccountType string    `json:"account_type"`
	LocalShare  *float64  `json:"local_share"`
	Notes       *string   `json:"notes"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts Account to AccountResponse
func (a *Account) ToResponse() *AccountResponse {
	return &AccountResponse{
		ID:          a.ID,
		AccountName: a.AccountName,
		AccountType: a.AccountType,
		LocalShare:  a.LocalShare,
		Notes:       a.Notes,
		IsActive:    a.IsActive,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
