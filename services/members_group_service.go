package services

import (
	"errors"
	"storeHouse/models"
	"storeHouse/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type MembersGroupService struct {
	DB *sqlx.DB
}

// Create a new instance of MembersGroupService
func NewMembersGroupService(db *sqlx.DB) *MembersGroupService {
	return &MembersGroupService{DB: db}
}

// CreateGroup handles group creation business logic
func (s *MembersGroupService) CreateGroup(req models.CreateGroupRequest, createdBy string) (*models.GroupResponse, error) {
	// Check for duplicate group name
	if _, err := repository.GetGroupByName(s.DB, req.GroupName); err == nil {
		return nil, errors.New("group name already exists")
	}

	// Prepare model for DB
	group := models.MembersGroup{
		GroupName: req.GroupName,
		Notes:     req.Notes,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to DB
	newGroup, err := repository.CreateGroup(s.DB, group)
	if err != nil {
		return nil, err
	}

	return newGroup.ToResponse(), nil
}

// UpdateGroup handles update logic
func (s *MembersGroupService) UpdateGroup(id string, req models.UpdateGroupRequest) (*models.GroupResponse, error) {
	// Fetch existing record
	existing, err := repository.GetGroup(s.DB, id)
	if err != nil {
		return nil, errors.New("group not found")
	}

	// Apply updates only if fields are provided
	if req.GroupName != nil {
		// Check for duplicate group name if different
		if existing.GroupName != *req.GroupName {
			if _, err := repository.GetGroupByName(s.DB, *req.GroupName); err == nil {
				return nil, errors.New("group name already exists")
			}
		}
		existing.GroupName = *req.GroupName
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	existing.UpdatedAt = time.Now()

	// Persist update
	updated, err := repository.UpdateGroup(s.DB, existing)
	if err != nil {
		return nil, err
	}

	return updated.ToResponse(), nil
}

// DeleteGroup removes a group record
func (s *MembersGroupService) DeleteGroup(id string) error {
	// Ensure exists before deleting
	if _, err := repository.GetGroup(s.DB, id); err != nil {
		return errors.New("group not found")
	}

	// Check if there are members in this group
	members, err := repository.GetMemberByGroup(s.DB, id)
	if err == nil && len(members) > 0 {
		return errors.New("cannot delete group with existing members")
	}

	return repository.DeleteGroup(s.DB, id)
}

// GetGroup returns single group details
func (s *MembersGroupService) GetGroup(id string) (*models.GroupResponse, error) {
	group, err := repository.GetGroup(s.DB, id)
	if err != nil {
		return nil, errors.New("group not found")
	}
	return group.ToResponse(), nil
}

// GetAllGroups returns all groups
func (s *MembersGroupService) GetAllGroups() ([]models.GroupResponse, error) {
	groups, err := repository.GetAllGroups(s.DB)
	if err != nil {
		return nil, err
	}

	// Convert to response list
	responses := make([]models.GroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, *g.ToResponse())
	}

	return responses, nil
}

// GetGroupByName returns group by name
func (s *MembersGroupService) GetGroupByName(name string) (*models.GroupResponse, error) {
	group, err := repository.GetGroupByName(s.DB, name)
	if err != nil {
		return nil, errors.New("group not found")
	}
	return group.ToResponse(), nil
}

// GetGroupsWithMemberCount returns all groups with their member counts
func (s *MembersGroupService) GetGroupsWithMemberCount() ([]map[string]interface{}, error) {
	return repository.GetGroupsWithMemberCount(s.DB)
}

// GetGroupMemberCount returns the number of members in a specific group
func (s *MembersGroupService) GetGroupMemberCount(groupID string) (int, error) {
	members, err := repository.GetMemberByGroup(s.DB, groupID)
	if err != nil {
		return 0, err
	}
	return len(members), nil
}
