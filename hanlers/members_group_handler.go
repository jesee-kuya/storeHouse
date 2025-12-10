package handlers

import (
	"encoding/json"
	"net/http"
	"storeHouse/models"
	"storeHouse/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type MembersGroupHandler struct {
	groupService *services.MembersGroupService
}

func NewMembersGroupHandler(db *sqlx.DB) *MembersGroupHandler {
	return &MembersGroupHandler{
		groupService: services.NewMembersGroupService(db),
	}
}

// CreateGroup handles group creation
func (h *MembersGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// For now, using a default user ID - in real app, get from authentication
	createdBy := "system"

	group, err := h.groupService.CreateGroup(req, createdBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

// GetGroup handles getting group by ID
func (h *MembersGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	group, err := h.groupService.GetGroup(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

// GetGroupByName handles getting group by name
func (h *MembersGroupHandler) GetGroupByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	group, err := h.groupService.GetGroupByName(name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

// GetAllGroups handles getting all groups
func (h *MembersGroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupService.GetAllGroups()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// GetGroupsWithMemberCount handles getting all groups with their member counts
func (h *MembersGroupHandler) GetGroupsWithMemberCount(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupService.GetGroupsWithMemberCount()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// GetGroupMemberCount handles getting the member count for a specific group
func (h *MembersGroupHandler) GetGroupMemberCount(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	count, err := h.groupService.GetGroupMemberCount(groupID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"member_count": count})
}

// UpdateGroup handles updating group details
func (h *MembersGroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.UpdateGroupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	group, err := h.groupService.UpdateGroup(id, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "group not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

// DeleteGroup handles deleting a group
func (h *MembersGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.groupService.DeleteGroup(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "group not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Group deleted successfully"})
}
