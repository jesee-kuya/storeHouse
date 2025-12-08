package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExpenditureService struct {
	DB *sqlx.DB
}

// Create a new instance of ExpenditureService
func NewExpenditureService(db *sqlx.DB) *ExpenditureService {
	return &ExpenditureService{DB: db}
}

// CreateExpenditure handles expenditure creation business logic
func (s *ExpenditureService) CreateExpenditure(req models.CreateExpenditureRequest, createdBy string) (*models.ExpenditureResponse, error) {
	// Validate that amount is positive
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Check if bank account exists
	if _, err := repository.GetAccount(s.DB, req.BankAccountID); err != nil {
		return nil, errors.New("bank account not found")
	}

	// Check if transaction exists
	if _, err := repository.GetTransaction(s.DB, req.TransactionID); err != nil {
		return nil, errors.New("transaction not found")
	}

	// Prepare model for DB
	expenditure := models.Expenditure{
		TransactionID: req.TransactionID,
		Particulars:   req.Particulars,
		BankAccountID: req.BankAccountID,
		Amount:        req.Amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to DB
	newExpenditure, err := repository.CreateExpenditure(s.DB, expenditure)
	if err != nil {
		return nil, err
	}

	return newExpenditure.ToResponse(), nil
}

// UpdateExpenditure handles update logic
func (s *ExpenditureService) UpdateExpenditure(id string, req models.UpdateExpenditureRequest) (*models.ExpenditureResponse, error) {
	// Fetch existing record
	existing, err := repository.GetExpenditure(s.DB, id)
	if err != nil {
		return nil, errors.New("expenditure not found")
	}

	// Apply updates only if fields are provided
	if req.Particulars != nil {
		existing.Particulars = *req.Particulars
	}
	if req.BankAccountID != nil {
		// Check if new bank account exists
		if _, err := repository.GetAccount(s.DB, *req.BankAccountID); err != nil {
			return nil, errors.New("bank account not found")
		}
		existing.BankAccountID = *req.BankAccountID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, errors.New("amount must be greater than zero")
		}
		existing.Amount = *req.Amount
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateExpenditure(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteExpenditure removes an expenditure record
func (s *ExpenditureService) DeleteExpenditure(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetExpenditure(s.DB, id); err != nil {
		return errors.New("expenditure not found")
	}

	return repository.DeleteExpenditure(s.DB, id)
}

// GetExpenditure returns single expenditure details
func (s *ExpenditureService) GetExpenditure(id string) (*models.ExpenditureResponse, error) {
	expenditure, err := repository.GetExpenditure(s.DB, id)
	if err != nil {
		return nil, errors.New("expenditure not found")
	}
	return expenditure.ToResponse(), nil
}

// GetAllExpenditures returns all expenditures
func (s *ExpenditureService) GetAllExpenditures() ([]models.ExpenditureResponse, error) {
	expenditures, err := repository.GetAllExpenditures(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ExpenditureResponse, 0, len(expenditures))
	for _, e := range expenditures {
		responses = append(responses, *e.ToResponse())
	}

	return responses, nil
}

// GetExpendituresByTransaction returns expenditures for a specific transaction
func (s *ExpenditureService) GetExpendituresByTransaction(transactionID string) ([]models.ExpenditureResponse, error) {
	expenditures, err := repository.GetExpenditureByTransaction(s.DB, transactionID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.ExpenditureResponse, 0, len(expenditures))
	for _, e := range expenditures {
		responses = append(responses, *e.ToResponse())
	}

	return responses, nil
}
