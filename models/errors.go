package models

import "errors"

// Common errors for the models package
var (
	ErrInvalidRole        = errors.New("invalid user role")
	ErrUserNotFound       = errors.New("user not found")
	ErrMemberNotFound     = errors.New("member not found")
	ErrGroupNotFound      = errors.New("group not found")
	ErrAccountNotFound    = errors.New("account not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrReceiptNotFound    = errors.New("receipt not found")
	ErrExpenditureNotFound = errors.New("expenditure not found")
	ErrTransferNotFound   = errors.New("transfer not found")
	
	// Validation errors
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPhone       = errors.New("invalid phone number format")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidAccountType = errors.New("invalid account type")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)