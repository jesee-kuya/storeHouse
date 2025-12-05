package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeGroupQuery(db *sqlx.DB, query string, group models.MembersGroup) (models.MembersGroup, error) {
	_, err := db.NamedExec(query, group)
	if err != nil {
		return models.MembersGroup{}, err
	}

	return group, nil
}

func CreateGroup(db *sqlx.DB, group models.MembersGroup) (models.MembersGroup, error) {
	group.ID = uuid.New().String()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	query := `INSERT INTO members_groups (id, group_name, notes, created_by, created_at, updated_at)
              VALUES (:id, :group_name, :notes, :created_by, :created_at, :updated_at)`

	return executeGroupQuery(db, query, group)
}

func UpdateGroup(db *sqlx.DB, group models.MembersGroup) (models.MembersGroup, error) {
	group.UpdatedAt = time.Now()

	query := `UPDATE members_groups SET group_name = :group_name, notes = :notes, updated_at = :updated_at 
			  WHERE id = :id`

	return executeGroupQuery(db, query, group)
}

func DeleteGroup(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM members_groups WHERE id = $1", id)
	return err
}

func GetGroup(db *sqlx.DB, id string) (models.MembersGroup, error) {
	var group models.MembersGroup
	err := db.Get(&group, "SELECT * FROM members_groups WHERE id = $1", id)
	if err != nil {
		return models.MembersGroup{}, err
	}

	return group, nil
}

func GetGroupByName(db *sqlx.DB, name string) (models.MembersGroup, error) {
	var group models.MembersGroup
	err := db.Get(&group, "SELECT * FROM members_groups WHERE group_name = $1", name)
	if err != nil {
		return models.MembersGroup{}, err
	}

	return group, nil
}

func GetAllGroups(db *sqlx.DB) ([]models.MembersGroup, error) {
	var groups []models.MembersGroup
	err := db.Select(&groups, "SELECT * FROM members_groups ORDER BY group_name ASC")
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetGroupsWithMemberCount(db *sqlx.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT 
			mg.id,
			mg.group_name,
			mg.notes,
			mg.created_by,
			mg.created_at,
			mg.updated_at,
			COUNT(m.id) as member_count
		FROM members_groups mg
		LEFT JOIN members m ON mg.id = m.group
		GROUP BY mg.id, mg.group_name, mg.notes, mg.created_by, mg.created_at, mg.updated_at
		ORDER BY mg.group_name ASC
	`
	err := db.Select(&results, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}
