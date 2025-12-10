package handlers

import (
	"encoding/json"
	"net/http"
	"storeHouse/models"
	"storeHouse/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type MemberHandler struct {
	memberService *services.MemberService
}

func NewMemberHandler(db *sqlx.DB) *MemberHandler {
	return &MemberHandler{
		memberService: services.NewMemberService(db),
	}
}

// CreateMember handles member creation
func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// For now, using a default user ID - in real app, get from authentication
	createdBy := "system"

	member, err := h.memberService.CreateMember(req, createdBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// GetMember handles getting member by ID
func (h *MemberHandler) GetMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	member, err := h.memberService.GetMember(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// GetMemberByPhone handles getting member by phone number
func (h *MemberHandler) GetMemberByPhone(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phone")

	member, err := h.memberService.GetMemberByPhone(phoneNumber)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// GetMemberByEmail handles getting member by email
func (h *MemberHandler) GetMemberByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	member, err := h.memberService.GetMemberByEmail(email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// GetAllMembers handles getting all members
func (h *MemberHandler) GetAllMembers(w http.ResponseWriter, r *http.Request) {
	members, err := h.memberService.GetAllMembers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// UpdateMember handles updating member details
func (h *MemberHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.UpdateMemberRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	member, err := h.memberService.UpdateMember(id, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "member not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// DeleteMember handles deleting a member
func (h *MemberHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.memberService.DeleteMember(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Member deleted successfully"})
}

// GetMembersByGroup handles getting members for a specific group
func (h *MemberHandler) GetMembersByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")

	members, err := h.memberService.GetMembersByGroup(groupID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// SearchMembers handles searching for members
func (h *MemberHandler) SearchMembers(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")

	if searchTerm == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "search term 'q' query parameter is required"})
		return
	}

	members, err := h.memberService.SearchMembers(searchTerm)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
