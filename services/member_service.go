package services

import (
	"errors"
	"regexp"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type MemberService struct {
	DB *sqlx.DB
}

// Create a new instance of MemberService
func NewMemberService(db *sqlx.DB) *MemberService {
	return &MemberService{DB: db}
}

// CreateMember handles member creation business logic
func (s *MemberService) CreateMember(req models.CreateMemberRequest, createdBy string) (*models.MemberResponse, error) {
	// Validate phone number format (basic validation)
	if !isValidPhoneNumber(req.PhoneNumber) {
		return nil, errors.New("invalid phone number format")
	}

	// Check for duplicate phone number
	if _, err := repository.GetMemberByPhone(s.DB, req.PhoneNumber); err == nil {
		return nil, errors.New("member with this phone number already exists")
	}

	// Check for duplicate email if provided
	if req.Email != nil {
		if _, err := repository.GetMemberByEmail(s.DB, *req.Email); err == nil {
			return nil, errors.New("member with this email already exists")
		}
	}

	// Validate group if provided
	if req.GroupID != nil {
		if _, err := repository.GetGroup(s.DB, *req.GroupID); err != nil {
			return nil, errors.New("specified group not found")
		}
	}

	// Prepare model for DB
	member := models.Member{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Notes:       req.Notes,
		GroupID:     req.GroupID,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to DB
	newMember, err := repository.CreateMember(s.DB, member)
	if err != nil {
		return nil, err
	}

	return newMember.ToResponse(), nil
}

// UpdateMember handles update logic
func (s *MemberService) UpdateMember(id string, req models.UpdateMemberRequest) (*models.MemberResponse, error) {
	// Fetch existing record
	existing, err := repository.GetMember(s.DB, id)
	if err != nil {
		return nil, errors.New("member not found")
	}

	// Apply updates only if fields are provided
	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.PhoneNumber != nil {
		// Validate phone number format
		if !isValidPhoneNumber(*req.PhoneNumber) {
			return nil, errors.New("invalid phone number format")
		}
		// Check for duplicate phone number
		if existing.PhoneNumber != *req.PhoneNumber {
			if _, err := repository.GetMemberByPhone(s.DB, *req.PhoneNumber); err == nil {
				return nil, errors.New("member with this phone number already exists")
			}
		}
		existing.PhoneNumber = *req.PhoneNumber
	}
	if req.Email != nil {
		// Check for duplicate email if different
		if existing.Email == nil || *existing.Email != *req.Email {
			if _, err := repository.GetMemberByEmail(s.DB, *req.Email); err == nil {
				return nil, errors.New("member with this email already exists")
			}
		}
		existing.Email = req.Email
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.GroupID != nil {
		// Validate group if provided
		if *req.GroupID != "" {
			if _, err := repository.GetGroup(s.DB, *req.GroupID); err != nil {
				return nil, errors.New("specified group not found")
			}
		}
		existing.GroupID = req.GroupID
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateMember(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteMember removes a member record
func (s *MemberService) DeleteMember(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetMember(s.DB, id); err != nil {
		return errors.New("member not found")
	}

	return repository.DeleteMember(s.DB, id)
}

// GetMember returns single member details
func (s *MemberService) GetMember(id string) (*models.MemberResponse, error) {
	member, err := repository.GetMember(s.DB, id)
	if err != nil {
		return nil, errors.New("member not found")
	}
	return member.ToResponse(), nil
}

// GetAllMembers returns all members
func (s *MemberService) GetAllMembers() ([]models.MemberResponse, error) {
	members, err := repository.GetAllMembers(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.MemberResponse, 0, len(members))
	for _, m := range members {
		responses = append(responses, *m.ToResponse())
	}

	return responses, nil
}

// GetMembersByGroup returns members in a specific group
func (s *MemberService) GetMembersByGroup(groupID string) ([]models.MemberResponse, error) {
	members, err := repository.GetMemberByGroup(s.DB, groupID)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.MemberResponse, 0, len(members))
	for _, m := range members {
		responses = append(responses, *m.ToResponse())
	}

	return responses, nil
}

// SearchMembers searches for members by name, phone number, or email
func (s *MemberService) SearchMembers(searchTerm string) ([]models.MemberResponse, error) {
	members, err := repository.SearchMembers(s.DB, searchTerm)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.MemberResponse, 0, len(members))
	for _, m := range members {
		responses = append(responses, *m.ToResponse())
	}

	return responses, nil
}

// GetMemberByPhone returns member by phone number
func (s *MemberService) GetMemberByPhone(phoneNumber string) (*models.MemberResponse, error) {
	member, err := repository.GetMemberByPhone(s.DB, phoneNumber)
	if err != nil {
		return nil, errors.New("member not found")
	}
	return member.ToResponse(), nil
}

// GetMemberByEmail returns member by email
func (s *MemberService) GetMemberByEmail(email string) (*models.MemberResponse, error) {
	member, err := repository.GetMemberByEmail(s.DB, email)
	if err != nil {
		return nil, errors.New("member not found")
	}
	return member.ToResponse(), nil
}

// isValidPhoneNumber performs basic phone number validation
func isValidPhoneNumber(phone string) bool {
	// Basic validation: allows digits, spaces, hyphens, parentheses, and plus sign
	phoneRegex := regexp.MustCompile(`^[\+]?[\d\s\-\(\)]{7,20}$`)
	return phoneRegex.MatchString(phone)
}
