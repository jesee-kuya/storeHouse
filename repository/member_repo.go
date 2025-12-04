package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeMemberQuery(db *sqlx.DB, query string, member models.Member) (models.Member, error) {
	_, err := db.NamedExec(query, member)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func CreateMember(db *sqlx.DB, member models.Member) (models.Member, error) {
	member.ID = uuid.New().String()
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()

	query := `INSERT INTO members (id, full_name, phone_number, email, notes, group, created_by, created_at, updated_at)
              VALUES (:id, :full_name, :phone_number, :email, :notes, :group_id, :created_by, :created_at, :updated_at)`

	return executeMemberQuery(db, query, member)
}

func UpdateMember(db *sqlx.DB, member models.Member) (models.Member, error) {
	member.UpdatedAt = time.Now()

	query := `UPDATE members SET full_name = :full_name, phone_number = :phone_number, email = :email, notes = :notes, group = :group_id, updated_at = :updated_at 
			  WHERE id = :id`

	return executeMemberQuery(db, query, member)
}

func DeleteMember(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM members WHERE id = $1", id)
	return err
}

func GetMember(db *sqlx.DB, id string) (models.Member, error) {
	var member models.Member
	err := db.Get(&member, "SELECT * FROM members WHERE id = $1", id)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func GetMemberByPhone(db *sqlx.DB, phoneNumber string) (models.Member, error) {
	var member models.Member
	err := db.Get(&member, "SELECT * FROM members WHERE phone_number = $1", phoneNumber)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func GetMemberByEmail(db *sqlx.DB, email string) (models.Member, error) {
	var member models.Member
	err := db.Get(&member, "SELECT * FROM members WHERE email = $1", email)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}

func GetMemberByGroup(db *sqlx.DB, groupID string) ([]models.Member, error) {
	var members []models.Member
	err := db.Select(&members, "SELECT * FROM members WHERE group = $1 ORDER BY full_name ASC", groupID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func GetAllMembers(db *sqlx.DB) ([]models.Member, error) {
	var members []models.Member
	err := db.Select(&members, "SELECT * FROM members ORDER BY full_name ASC")
	if err != nil {
		return nil, err
	}

	return members, nil
}

func SearchMembers(db *sqlx.DB, searchTerm string) ([]models.Member, error) {
	var members []models.Member
	query := "SELECT * FROM members WHERE full_name ILIKE $1 OR phone_number ILIKE $1 OR email ILIKE $1 ORDER BY full_name ASC"
	err := db.Select(&members, query, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}

	return members, nil
}
