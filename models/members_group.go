package models

import (
	"time"
)

// MembersGroup represents church member groups or estates
type MembersGroup struct {
	ID          string    `json:"id" db:"id"`
	GroupName   string    `json:"group_name" db:"group_name" binding:"required,max=50"`
	Notes       *string   `json:"notes" db:"notes"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateGroupRequest represents the request for creating a new group
type CreateGroupRequest struct {
	GroupName string  `json:"group_name" binding:"required,max=50"`
	Notes     *string `json:"notes"`
}

// UpdateGroupRequest represents the request for updating a group
type UpdateGroupRequest struct {
	GroupName *string `json:"group_name" binding:"max=50"`
	Notes     *string `json:"notes"`
}

// GroupResponse represents the group response
type GroupResponse struct {
	ID        string    `json:"id"`
	GroupName string    `json:"group_name"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts MembersGroup to GroupResponse
func (g *MembersGroup) ToResponse() *GroupResponse {
	return &GroupResponse{
		ID:        g.ID,
		GroupName: g.GroupName,
		Notes:     g.Notes,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}