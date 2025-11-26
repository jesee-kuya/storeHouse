package models

import (
	"time"
)

// Expenditure represents expenses and withdrawals from accounts
type Expenditure struct {
	ID            string      `json:"id" db:"id"`
	TransactionID string      `json:"transaction_id" db:"transaction_id" binding:"required"`
	Transaction   *Transaction `json:"transaction,omitempty" db:"-"`
	Particulars   string      `json:"particulars" db:"perticulars" binding:"required,max=255"`
	BankAccountID string      `json:"bank_account_id" db:"bank_account" binding:"required"`
	BankAccount   *Account    `json:"bank_account,omitempty" db:"-"`
	Amount        float64     `json:"amount" db:"amount" binding:"required"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateExpenditureRequest represents the request for creating a new expenditure
type CreateExpenditureRequest struct {
	TransactionID string  `json:"transaction_id" binding:"required"`
	Particulars   string  `json:"particulars" binding:"required,max=255"`
	BankAccountID string  `json:"bank_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

// UpdateExpenditureRequest represents the request for updating an expenditure
type UpdateExpenditureRequest struct {
	Particulars   *string  `json:"particulars" binding:"max=255"`
	BankAccountID *string  `json:"bank_account_id"`
	Amount        *float64 `json:"amount"`
}

// ExpenditureResponse represents the expenditure response
type ExpenditureResponse struct {
	ID            string             `json:"id"`
	TransactionID string             `json:"transaction_id"`
	Transaction   *TransactionResponse `json:"transaction,omitempty"`
	Particulars   string             `json:"particulars"`
	BankAccountID string             `json:"bank_account_id"`
	BankAccount   *AccountResponse   `json:"bank_account,omitempty"`
	Amount        float64            `json:"amount"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// ToResponse converts Expenditure to ExpenditureResponse
func (e *Expenditure) ToResponse() *ExpenditureResponse {
	var transactionResp *TransactionResponse
	if e.Transaction != nil {
		transactionResp = e.Transaction.ToResponse()
	}
	
	var bankAccountResp *AccountResponse
	if e.BankAccount != nil {
		bankAccountResp = e.BankAccount.ToResponse()
	}
	
	return &ExpenditureResponse{
		ID:             e.ID,
		TransactionID:  e.TransactionID,
		Transaction:    transactionResp,
		Particulars:    e.Particulars,
		BankAccountID:  e.BankAccountID,
		BankAccount:    bankAccountResp,
		Amount:         e.Amount,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}