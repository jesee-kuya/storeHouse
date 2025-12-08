package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransferService struct {
	DB *sqlx.DB
}

// Create a new instance of TransferService
func NewTransferService(db *sqlx.DB) *TransferService {
	return &TransferService{DB: db}
}

// CreateTransfer handles transfer creation business logic
func (s *TransferService) CreateTransfer(req models.CreateTransferRequest) (*models.TransferResponse, error) {
	// Validate that amount is positive
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Check if credit account exists
	if _, err := repository.GetAccount(s.DB, req.CreditAccountID); err != nil {
		return nil, errors.New("credit account not found")
	}

	// Check if transaction exists
	if _, err := repository.GetTransaction(s.DB, req.TransactionID); err != nil {
		return nil, errors.New("transaction not found")
	}

	// Prepare model for DB
	transfer := models.Transfer{
		TransactionID:   req.TransactionID,
		Particulars:     req.Particulars,
		CreditAccountID: req.CreditAccountID,
		Amount:          req.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to DB
	newTransfer, err := repository.CreateTransfer(s.DB, transfer)
	if err != nil {
		return nil, err
	}

	return newTransfer.ToResponse(), nil
}

// UpdateTransfer handles update logic
func (s *TransferService) UpdateTransfer(id string, req models.UpdateTransferRequest) (*models.TransferResponse, error) {
	// Fetch existing record
	existing, err := repository.GetTransfer(s.DB, id)
	if err != nil {
		return nil, errors.New("transfer not found")
	}

	// Apply updates only if fields are provided
	if req.Particulars != nil {
		existing.Particulars = *req.Particulars
	}
	if req.CreditAccountID != nil {
		// Check if new credit account exists
		if _, err := repository.GetAccount(s.DB, *req.CreditAccountID); err != nil {
			return nil, errors.New("credit account not found")
		}
		existing.CreditAccountID = *req.CreditAccountID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, errors.New("amount must be greater than zero")
		}
		existing.Amount = *req.Amount
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateTransfer(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteTransfer removes a transfer record
func (s *TransferService) DeleteTransfer(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetTransfer(s.DB, id); err != nil {
		return errors.New("transfer not found")
	}

	return repository.DeleteTransfer(s.DB, id)
}

// GetTransfer returns single transfer details
func (s *TransferService) GetTransfer(id string) (*models.TransferResponse, error) {
	transfer, err := repository.GetTransfer(s.DB, id)
	if err != nil {
		return nil, errors.New("transfer not found")
	}
	return transfer.ToResponse(), nil
}

// GetAllTransfers returns all transfers
func (s *TransferService) GetAllTransfers() ([]models.TransferResponse, error) {
	transfers, err := repository.GetAllTransfers(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransferResponse, 0, len(transfers))
	for _, t := range transfers {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransfersByTransaction returns transfers for a specific transaction
func (s *TransferService) GetTransfersByTransaction(transactionID string) ([]models.TransferResponse, error) {
	transfers, err := repository.GetTransferByTransaction(s.DB, transactionID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransferResponse, 0, len(transfers))
	for _, t := range transfers {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTransfersByCreditAccount returns transfers for a specific credit account
func (s *TransferService) GetTransfersByCreditAccount(accountID string) ([]models.TransferResponse, error) {
	transfers, err := repository.GetTransferByCreditAccount(s.DB, accountID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransferResponse, 0, len(transfers))
	for _, t := range transfers {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTotalTransfersByCreditAccount returns the total amount for a specific credit account
func (s *TransferService) GetTotalTransfersByCreditAccount(accountID string) (float64, error) {
	return repository.GetTotalTransfersByCreditAccount(s.DB, accountID)
}

// GetTransfersByDateRange returns transfers within a date range
func (s *TransferService) GetTransfersByDateRange(startDate, endDate time.Time) ([]models.TransferResponse, error) {
	transfers, err := repository.GetTransferByDateRange(s.DB, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.TransferResponse, 0, len(transfers))
	for _, t := range transfers {
		responses = append(responses, *t.ToResponse())
	}

	return responses, nil
}

// GetTotalTransfersByDateRange returns the total transfers for an account within a date range
func (s *TransferService) GetTotalTransfersByDateRange(accountID string, startDate, endDate time.Time) (float64, error) {
	return repository.GetTotalTransfersByDateRange(s.DB, accountID, startDate, endDate)
}