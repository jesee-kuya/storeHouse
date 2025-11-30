package models

import (
	"time"
)

// Receipt represents income/receipts from offerings
type Receipt struct {
	ID            string     `json:"id" db:"id"`
	TransactionID string     `json:"transaction_id" db:"transaction_id" binding:"required"`
	Transaction   *Transaction `json:"transaction,omitempty" db:"-"`
	IncomeAccountID string   `json:"income_account_id" db:"income_account" binding:"required"`
	IncomeAccount *Account   `json:"income_account,omitempty" db:"-"`
	Amount        float64    `json:"amount" db:"amount" binding:"required"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateReceiptRequest represents the request for creating a new receipt
type CreateReceiptRequest struct {
	TransactionID string  `json:"transaction_id" binding:"required"`
	IncomeAccountID string `json:"income_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

// UpdateReceiptRequest represents the request for updating a receipt
type UpdateReceiptRequest struct {
	IncomeAccountID *string  `json:"income_account_id"`
	Amount          *float64 `json:"amount"`
}

// ReceiptResponse represents the receipt response
type ReceiptResponse struct {
	ID            string          `json:"id"`
	TransactionID string          `json:"transaction_id"`
	Transaction   *TransactionResponse `json:"transaction,omitempty"`
	IncomeAccountID string        `json:"income_account_id"`
	IncomeAccount *AccountResponse `json:"income_account,omitempty"`
	Amount        float64         `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ToResponse converts Receipt to ReceiptResponse
func (r *Receipt) ToResponse() *ReceiptResponse {
	var transactionResp *TransactionResponse
	if r.Transaction != nil {
		transactionResp = r.Transaction.ToResponse()
	}
	
	var incomeAccountResp *AccountResponse
	if r.IncomeAccount != nil {
		incomeAccountResp = r.IncomeAccount.ToResponse()
	}
	
	return &ReceiptResponse{
		ID:               r.ID,
		TransactionID:    r.TransactionID,
		Transaction:      transactionResp,
		IncomeAccountID:  r.IncomeAccountID,
		IncomeAccount:    incomeAccountResp,
		Amount:           r.Amount,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}