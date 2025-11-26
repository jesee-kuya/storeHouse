package models

import (
	"time"
)

// Transaction represents main transaction records
type Transaction struct {
	ID              string        `json:"id" db:"id"`
	TransactionRef  *string       `json:"transaction_ref" db:"transaction_ref" binding:"max=20"`
	TransactionDate time.Time     `json:"transaction_date" db:"transaction_date"`
	TransactionType string        `json:"transaction_type" db:"transaction_type" binding:"required"`
	Amount          float64       `json:"amount" db:"amount" binding:"required"`
	Notes           *string       `json:"notes" db:"notes"`
	DebitAccountID  string        `json:"debit_account_id" db:"debit_account" binding:"required"`
	DebitAccount    *Account      `json:"debit_account,omitempty" db:"-"`
	MemberID        *string       `json:"member_id" db:"member"`
	Member          *Member       `json:"member,omitempty" db:"-"`
	CreatedBy       string        `json:"created_by" db:"created_by" binding:"required"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

// TransactionType represents the different transaction types
type TransactionType string

const (
	TransactionReceipts    TransactionType = "receipts"
	TransactionWithdrawal TransactionType = "withdrawal"
	TransactionExpenses   TransactionType = "expenses"
	TransactionTransfer   TransactionType = "transfer"
)

// ValidateTransactionType checks if the transaction type is valid
func (t *Transaction) ValidateTransactionType() error {
	switch t.TransactionType {
	case string(TransactionReceipts), string(TransactionWithdrawal), string(TransactionExpenses), string(TransactionTransfer):
		return nil
	default:
		return ErrInvalidTransactionType
	}
}

// CreateTransactionRequest represents the request for creating a new transaction
type CreateTransactionRequest struct {
	TransactionRef  *string  `json:"transaction_ref" binding:"max=20"`
	TransactionDate *time.Time `json:"transaction_date"`
	TransactionType string    `json:"transaction_type" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Notes           *string   `json:"notes"`
	DebitAccountID  string    `json:"debit_account_id" binding:"required"`
	MemberID        *string   `json:"member_id"`
}

// UpdateTransactionRequest represents the request for updating a transaction
type UpdateTransactionRequest struct {
	TransactionRef  *string  `json:"transaction_ref" binding:"max=20"`
	TransactionDate *time.Time `json:"transaction_date"`
	TransactionType *string  `json:"transaction_type"`
	Amount          *float64 `json:"amount"`
	Notes           *string  `json:"notes"`
	DebitAccountID  *string  `json:"debit_account_id"`
	MemberID        *string  `json:"member_id"`
}

// TransactionResponse represents the transaction response
type TransactionResponse struct {
	ID              string          `json:"id"`
	TransactionRef  *string         `json:"transaction_ref"`
	TransactionDate time.Time       `json:"transaction_date"`
	TransactionType string          `json:"transaction_type"`
	Amount          float64         `json:"amount"`
	Notes           *string         `json:"notes"`
	DebitAccountID  string          `json:"debit_account_id"`
	DebitAccount    *AccountResponse `json:"debit_account,omitempty"`
	MemberID        *string         `json:"member_id"`
	Member          *MemberResponse  `json:"member,omitempty"`
	CreatedBy       string          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// ToResponse converts Transaction to TransactionResponse
func (t *Transaction) ToResponse() *TransactionResponse {
	var debitAccountResp *AccountResponse
	if t.DebitAccount != nil {
		debitAccountResp = t.DebitAccount.ToResponse()
	}
	
	var memberResp *MemberResponse
	if t.Member != nil {
		memberResp = t.Member.ToResponse()
	}
	
	return &TransactionResponse{
		ID:              t.ID,
		TransactionRef:  t.TransactionRef,
		TransactionDate: t.TransactionDate,
		TransactionType: t.TransactionType,
		Amount:          t.Amount,
		Notes:           t.Notes,
		DebitAccountID:  t.DebitAccountID,
		DebitAccount:    debitAccountResp,
		MemberID:        t.MemberID,
		Member:          memberResp,
		CreatedBy:       t.CreatedBy,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}