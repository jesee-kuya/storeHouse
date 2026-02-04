package models

import (
	"time"
)

// Member represents a church member who makes offerings and contributions
type Member struct {
	ID          string        `json:"id" db:"id"`
	FullName    string        `json:"full_name" db:"full_name" binding:"required,max=100"`
	PhoneNumber string        `json:"phone_number" db:"phone_number" binding:"required,max=20"`
	Email       *string       `json:"email" db:"email"`
	Notes       *string       `json:"notes" db:"notes"`
	GroupID     *string       `json:"group_id" db:"group_id"`
	Group       *MembersGroup `json:"group,omitempty" db:"-"`
	CreatedBy   string        `json:"created_by" db:"created_by"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// CreateMemberRequest represents the request for creating a new member
type CreateMemberRequest struct {
	FullName    string  `json:"full_name" binding:"required,max=100"`
	PhoneNumber string  `json:"phone_number" binding:"required,max=20"`
	Email       *string `json:"email" binding:"email"`
	Notes       *string `json:"notes"`
	GroupID     *string `json:"group_id"`
}

// UpdateMemberRequest represents the request for updating a member
type UpdateMemberRequest struct {
	FullName    *string `json:"full_name" binding:"max=100"`
	PhoneNumber *string `json:"phone_number" binding:"max=20"`
	Email       *string `json:"email" binding:"email"`
	Notes       *string `json:"notes"`
	GroupID     *string `json:"group_id"`
}

// MemberResponse represents the member response
type MemberResponse struct {
	ID          string         `json:"id"`
	FullName    string         `json:"full_name"`
	PhoneNumber string         `json:"phone_number"`
	Email       *string        `json:"email"`
	Notes       *string        `json:"notes"`
	GroupID     *string        `json:"group_id"`
	Group       *GroupResponse `json:"group,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ToResponse converts Member to MemberResponse
func (m *Member) ToResponse() *MemberResponse {
	var groupResp *GroupResponse
	if m.Group != nil {
		groupResp = m.Group.ToResponse()
	}
	
	return &MemberResponse{
		ID:          m.ID,
		FullName:    m.FullName,
		PhoneNumber: m.PhoneNumber,
		Email:       m.Email,
		Notes:       m.Notes,
		GroupID:     m.GroupID,
		Group:       groupResp,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}