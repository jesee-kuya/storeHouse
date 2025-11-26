package models

import (
	"time"
)

// Transfer represents account transfers
type Transfer struct {
	ID            string      `json:"id" db:"id"`
	TransactionID string      `json:"transaction_id" db:"transaction_id" binding:"required"`
	Transaction   *Transaction `json:"transaction,omitempty" db:"-"`
	Particulars   string      `json:"particulars" db:"perticulars" binding:"required,max=255"`
	CreditAccountID string    `json:"credit_account_id" db:"credit_account" binding:"required"`
	CreditAccount *Account    `json:"credit_account,omitempty" db:"-"`
	Amount        float64     `json:"amount" db:"amount" binding:"required"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateTransferRequest represents the request for creating a new transfer
type CreateTransferRequest struct {
	TransactionID    string  `json:"transaction_id" binding:"required"`
	Particulars      string  `json:"particulars" binding:"required,max=255"`
	CreditAccountID  string  `json:"credit_account_id" binding:"required"`
	Amount           float64 `json:"amount" binding:"required"`
}

// UpdateTransferRequest represents the request for updating a transfer
type UpdateTransferRequest struct {
	Particulars      *string  `json:"particulars" binding:"max=255"`
	CreditAccountID  *string  `json:"credit_account_id"`
	Amount           *float64 `json:"amount"`
}

// TransferResponse represents the transfer response
type TransferResponse struct {
	ID               string             `json:"id"`
	TransactionID    string             `json:"transaction_id"`
	Transaction      *TransactionResponse `json:"transaction,omitempty"`
	Particulars      string             `json:"particulars"`
	CreditAccountID  string             `json:"credit_account_id"`
	CreditAccount    *AccountResponse   `json:"credit_account,omitempty"`
	Amount           float64            `json:"amount"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// ToResponse converts Transfer to TransferResponse
func (t *Transfer) ToResponse() *TransferResponse {
	var transactionResp *TransactionResponse
	if t.Transaction != nil {
		transactionResp = t.Transaction.ToResponse()
	}
	
	var creditAccountResp *AccountResponse
	if t.CreditAccount != nil {
		creditAccountResp = t.CreditAccount.ToResponse()
	}
	
	return &TransferResponse{
		ID:               t.ID,
		TransactionID:    t.TransactionID,
		Transaction:      transactionResp,
		Particulars:      t.Particulars,
		CreditAccountID:  t.CreditAccountID,
		CreditAccount:    creditAccountResp,
		Amount:           t.Amount,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}