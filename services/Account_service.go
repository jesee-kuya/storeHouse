package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccountService struct {
	DB *sqlx.DB
}

// Create a new instance of AccountService
func NewAccountService(db *sqlx.DB) *AccountService {
	return &AccountService{DB: db}
}

// Create Account handles account creation business logic
func (s *AccountService) CreateAccount(req models.CreateAccountRequest, createdBy string) (*models.AccountResponse, error) {
	// Validate account type
	if err := (&models.Account{AccountType: req.AccountType}).ValidateAccountType(); err != nil {
		return nil, err
	}

	// Prepare model for DB
	account := models.Account{
		AccountName: req.AccountName,
		AccountType: req.AccountType,
		LocalShare:  req.LocalShare,
		Notes:       req.Notes,
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Check for duplicates before creating
	if _, err := repository.GetAccountByName(s.DB, req.AccountName); err == nil {
		return nil, errors.New("account name already exists")
	}

	// Save to DB
	newAcc, err := repository.CreateAccount(s.DB, account)
	if err != nil {
		return nil, err
	}

	return newAcc.ToResponse(), nil
}

// UpdateAccount handles update logic
func (s *AccountService) UpdateAccount(id string, req models.UpdateAccountRequest) (*models.AccountResponse, error) {
	// Fetch existing record
	existing, err := repository.GetAccount(s.DB, id)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Apply updates only if fields are provided
	if req.AccountName != nil {
		existing.AccountName = *req.AccountName
	}
	if req.AccountType != nil {
		existing.AccountType = *req.AccountType
		if err := existing.ValidateAccountType(); err != nil {
			return nil, err
		}
	}
	if req.LocalShare != nil {
		existing.LocalShare = req.LocalShare
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateAccount(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeactivateAccount sets is_active = false
func (s *AccountService) DeactivateAccount(id string) error {
	// Ensure exists before disabling
	if _, err := repository.GetAccount(s.DB, id); err != nil {
		return errors.New("account not found")
	}

	_, err := repository.DeactivateAccount(s.DB, id)
	return err
}

// GetAccount returns single account details
func (s *AccountService) GetAccount(id string) (*models.AccountResponse, error) {
	acc, err := repository.GetAccount(s.DB, id)
	if err != nil {
		return nil, errors.New("account not found")
	}
	return acc.ToResponse(), nil
}

// List all active accounts
func (s *AccountService) GetAllAccounts() ([]models.AccountResponse, error) {
	accounts, err := repository.GetAllAccounts(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.AccountResponse, 0, len(accounts))
	for _, a := range accounts {
		responses = append(responses, *a.ToResponse())
	}

	return responses, nil
}

