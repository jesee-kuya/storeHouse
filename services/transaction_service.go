package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	DB *sqlx.DB
}

// Create a new instance of TransactionService
func NewTransactionService(db *sqlx.DB) *TransactionService {
	return &TransactionService{DB: db}
}

// CreateTransaction handles transaction creation business logic
func (s *TransactionService) CreateTransaction(req models.CreateTransactionRequest, createdBy string) (*models.TransactionResponse, error) {
	// Validate transaction type
	transaction := models.Transaction{
		TransactionType: req.TransactionType,
	}
	if err := transaction.ValidateTransactionType(); err != nil {
		return nil, err
	}

	// Validate that amount is positive
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Check if debit account exists
	if _, err := repository.GetAccount(s.DB, req.DebitAccountID); err != nil {
		return nil, errors.New("debit account not found")
	}

	// Validate member if provided
	if req.MemberID != nil {
		if _, err := repository.GetMember(s.DB, *req.MemberID); err != nil {
			return nil, errors.New("member not found")
		}
	}

	// Check for duplicate transaction reference if provided
	if req.TransactionRef != nil {
		if _, err := repository.GetTransactionByRef(s.DB, *req.TransactionRef); err == nil {
			return nil, errors.New("transaction reference already exists")
		}
	}

	// Prepare model for DB
	transactionModel := models.Transaction{
		TransactionRef:  req.TransactionRef,
		TransactionDate: time.Now(),
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		Notes:           req.Notes,
		DebitAccountID:  req.DebitAccountID,
		MemberID:        req.MemberID,
		CreatedBy:       createdBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Use provided transaction date if available
	if req.TransactionDate != nil && !req.TransactionDate.IsZero() {
		transactionModel.TransactionDate = *req.TransactionDate
	}

	// Save to DB
	newTransaction, err := repository.CreateTransaction(s.DB, transactionModel)
	if err != nil {
		return nil, err
	}

	return newTransaction.ToResponse(), nil
}

// UpdateTransaction handles update logic
func (s *TransactionService) UpdateTransaction(id string, req models.UpdateTransactionRequest) (*models.TransactionResponse, error) {
	// Fetch existing record
	existing, err := repository.GetTransaction(s.DB, id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Apply updates only if fields are provided
	if req.TransactionRef != nil {
		// Check for duplicate transaction reference if different
		if existing.TransactionRef == nil || *existing.TransactionRef != *req.TransactionRef {
			if _, err := repository.GetTransactionByRef(s.DB, *req.TransactionRef); err == nil {
				return nil, errors.New("transaction reference already exists")
			}
		}
		existing.TransactionRef = req.TransactionRef
	}
	if req.TransactionDate != nil && !req.TransactionDate.IsZero() {
		existing.TransactionDate = *req.TransactionDate
	}
	if req.TransactionType != nil {
		existing.TransactionType = *req.TransactionType
		// Validate transaction type
		if err := existing.ValidateTransactionType(); err != nil {
			return nil, err
		}
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, errors.New("amount must be greater than zero")
		}
		existing.Amount = *req.Amount
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.DebitAccountID != nil {
		// Check if new debit account exists
		if _, err := repository.GetAccount(s.DB, *req.DebitAccountID); err != nil {
			return nil, errors.New("debit account not found")
		}
		existing.DebitAccountID = *req.DebitAccountID
	}
	if req.MemberID != nil {
		// Validate member if provided
		if *req.MemberID != "" {
			if _, err := repository.GetMember(s.DB, *req.MemberID); err != nil {
				return nil, errors.New("member not found")
			}
		}
		existing.MemberID = req.MemberID
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateTransaction(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteTransaction removes a transaction record
func (s *TransactionService) DeleteTransaction(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetTransaction(s.DB, id); err != nil {
		return errors.New("transaction not found")
	}

	return repository.DeleteTransaction(s.DB, id)
}

// GetTransaction returns single transaction details
func (s *TransactionService) GetTransaction(id string) (*models.TransactionResponse, error) {
	transaction, err := repository.GetTransaction(s.DB, id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return transaction.ToResponse(), nil
}

// GetTransactionByRef returns transaction by reference
func (s *TransactionService) GetTransactionByRef(ref string) (*models.TransactionResponse, error) {
	transaction, err := repository.GetTransactionByRef(s.DB, ref)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return transaction.ToResponse(), nil
}

// GetAllTransactions returns all transactions
func (s *TransactionService) GetAllTransactions() ([]models.TransactionResponse, error) {
	transactions, err := repository.GetAllTransactions(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransactionsByAccount returns transactions for a specific debit account
func (s *TransactionService) GetTransactionsByAccount(accountID string) ([]models.TransactionResponse, error) {
	transactions, err := repository.GetTransactionsByAccount(s.DB, accountID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransactionsByMember returns transactions for a specific member
func (s *TransactionService) GetTransactionsByMember(memberID string) ([]models.TransactionResponse, error) {
	transactions, err := repository.GetTransactionsByMember(s.DB, memberID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransactionsByType returns transactions of a specific type
func (s *TransactionService) GetTransactionsByType(transactionType string) ([]models.TransactionResponse, error) {
	// Validate transaction type
	transaction := models.Transaction{TransactionType: transactionType}
	if err := transaction.ValidateTransactionType(); err != nil {
		return nil, err
	}

	transactions, err := repository.GetTransactionsByType(s.DB, transactionType)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransactionsByDateRange returns transactions within a date range
func (s *TransactionService) GetTransactionsByDateRange(startDate, endDate time.Time) ([]models.TransactionResponse, error) {
	transactions, err := repository.GetTransactionsByDateRange(s.DB, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}
