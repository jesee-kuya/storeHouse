package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReceiptService struct {
	DB *sqlx.DB
}

// Create a new instance of ReceiptService
func NewReceiptService(db *sqlx.DB) *ReceiptService {
	return &ReceiptService{DB: db}
}

// CreateReceipt handles receipt creation business logic
func (s *ReceiptService) CreateReceipt(req models.CreateReceiptRequest) (*models.ReceiptResponse, error) {
	// Validate that amount is positive
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Check if income account exists
	if _, err := repository.GetAccount(s.DB, req.IncomeAccountID); err != nil {
		return nil, errors.New("income account not found")
	}

	// Check if transaction exists
	if _, err := repository.GetTransaction(s.DB, req.TransactionID); err != nil {
		return nil, errors.New("transaction not found")
	}

	// Prepare model for DB
	receipt := models.Receipt{
		TransactionID:   req.TransactionID,
		IncomeAccountID: req.IncomeAccountID,
		Amount:          req.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to DB
	newReceipt, err := repository.CreateReceipt(s.DB, receipt)
	if err != nil {
		return nil, err
	}

	return newReceipt.ToResponse(), nil
}

// UpdateReceipt handles update logic
func (s *ReceiptService) UpdateReceipt(id string, req models.UpdateReceiptRequest) (*models.ReceiptResponse, error) {
	// Fetch existing record
	existing, err := repository.GetReceipt(s.DB, id)
	if err != nil {
		return nil, errors.New("receipt not found")
	}

	// Apply updates only if fields are provided
	if req.IncomeAccountID != nil {
		// Check if new income account exists
		if _, err := repository.GetAccount(s.DB, *req.IncomeAccountID); err != nil {
			return nil, errors.New("income account not found")
		}
		existing.IncomeAccountID = *req.IncomeAccountID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, errors.New("amount must be greater than zero")
		}
		existing.Amount = *req.Amount
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateReceipt(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteReceipt removes a receipt record
func (s *ReceiptService) DeleteReceipt(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetReceipt(s.DB, id); err != nil {
		return errors.New("receipt not found")
	}

	return repository.DeleteReceipt(s.DB, id)
}

// GetReceipt returns single receipt details
func (s *ReceiptService) GetReceipt(id string) (*models.ReceiptResponse, error) {
	receipt, err := repository.GetReceipt(s.DB, id)
	if err != nil {
		return nil, errors.New("receipt not found")
	}
	return receipt.ToResponse(), nil
}

// GetAllReceipts returns all receipts
func (s *ReceiptService) GetAllReceipts() ([]models.ReceiptResponse, error) {
	receipts, err := repository.GetAllReceipts(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ReceiptResponse, 0, len(receipts))
	for _, r := range receipts {
		responses = append(responses, *r.ToResponse())
	}

	return responses, nil
}

// GetReceiptsByTransaction returns receipts for a specific transaction
func (s *ReceiptService) GetReceiptsByTransaction(transactionID string) ([]models.ReceiptResponse, error) {
	receipts, err := repository.GetReceiptByTransaction(s.DB, transactionID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ReceiptResponse, 0, len(receipts))
	for _, r := range receipts {
		responses = append(responses, *r.ToResponse())
	}

	return responses, nil
}

// GetReceiptsByAccount returns receipts for a specific income account
func (s *ReceiptService) GetReceiptsByAccount(accountID string) ([]models.ReceiptResponse, error) {
	receipts, err := repository.GetReceiptByAccount(s.DB, accountID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ReceiptResponse, 0, len(receipts))
	for _, r := range receipts {
		responses = append(responses, *r.ToResponse())
	}

	return responses, nil
}

// GetTotalReceiptsByAccount returns the total amount for a specific account
func (s *ReceiptService) GetTotalReceiptsByAccount(accountID string) (float64, error) {
	return repository.GetTotalReceiptsByAccount(s.DB, accountID)
}

// GetReceiptsByDateRange returns receipts within a date range
func (s *ReceiptService) GetReceiptsByDateRange(startDate, endDate time.Time) ([]models.ReceiptResponse, error) {
	receipts, err := repository.GetReceiptByDateRange(s.DB, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ReceiptResponse, 0, len(receipts))
	for _, r := range receipts {
		responses = append(responses, *r.ToResponse())
	}

	return responses, nil
}

// GetTotalReceiptsByDateRange returns the total receipts for an account within a date range
func (s *ReceiptService) GetTotalReceiptsByDateRange(accountID string, startDate, endDate time.Time) (float64, error) {
	return repository.GetTotalReceiptsByDateRange(s.DB, accountID, startDate, endDate)
}
