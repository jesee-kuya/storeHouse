package handlers

import (
	"encoding/json"
	"net/http"
	"storeHouse/models"
	"storeHouse/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type ExpenditureHandler struct {
	expenditureService *services.ExpenditureService
}

func NewExpenditureHandler(db *sqlx.DB) *ExpenditureHandler {
	return &ExpenditureHandler{
		expenditureService: services.NewExpenditureService(db),
	}
}

// CreateExpenditure handles expenditure creation
func (h *ExpenditureHandler) CreateExpenditure(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpenditureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// For now, using a default user ID - in real app, get from authentication
	createdBy := "system"

	expenditure, err := h.expenditureService.CreateExpenditure(req, createdBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expenditure)
}

// GetExpenditure handles getting expenditure by ID
func (h *ExpenditureHandler) GetExpenditure(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	expenditure, err := h.expenditureService.GetExpenditure(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditure)
}

// GetAllExpenditures handles getting all expenditures
func (h *ExpenditureHandler) GetAllExpenditures(w http.ResponseWriter, r *http.Request) {
	expenditures, err := h.expenditureService.GetAllExpenditures()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditures)
}

// GetExpendituresByTransaction handles getting expenditures for a specific transaction
func (h *ExpenditureHandler) GetExpendituresByTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionID")

	expenditures, err := h.expenditureService.GetExpendituresByTransaction(transactionID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditures)
}

// UpdateExpenditure handles updating expenditure details
func (h *ExpenditureHandler) UpdateExpenditure(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.UpdateExpenditureRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	expenditure, err := h.expenditureService.UpdateExpenditure(id, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "expenditure not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenditure)
}

// DeleteExpenditure handles deleting an expenditure
func (h *ExpenditureHandler) DeleteExpenditure(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.expenditureService.DeleteExpenditure(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Expenditure deleted successfully"})
}
