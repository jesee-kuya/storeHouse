package handlers

import (
	"encoding/json"
	"net/http"
	"storeHouse/models"
	"storeHouse/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(db *sqlx.DB) *AccountHandler {
	return &AccountHandler{
		accountService: services.NewAccountService(db),
	}
}

// CreateAccount handles account creation
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	// For now, using a default user ID - in real app, get from authentication
	createdBy := "system"

	account, err := h.accountService.CreateAccount(req, createdBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetAccount handles getting account by ID
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	account, err := h.accountService.GetAccount(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// GetAllAccounts handles getting all accounts
func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// UpdateAccount handles updating account details
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.UpdateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.UpdateAccount(id, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "account not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// DeactivateAccount handles deactivating an account
func (h *AccountHandler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.accountService.DeactivateAccount(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Account deactivated successfully"})
}
