package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)



func exercuteQuery(db *sqlx.DB, query string, acc models.Account)  (models.Account, error) {
	_, err := db.NamedExec(query, acc)
	if err != nil {
		return models.Account{}, err
	}

	return acc, nil
}

func CreateAccount(db *sqlx.DB, acc models.Account) (models.Account, error) {
	acc.ID = uuid.New().String()
	acc.CreatedAt = time.Now()
	acc.UpdatedAt = time.Now()

	query := `INSERT INTO accounts (id, account_name, account_type, local_share, notes, is_active, created_at, updated_at)
              VALUES (:id, :account_name, :account_type, :local_share, :notes, :is_active, :created_at, :updated_at)`
	
	return exercuteQuery(db, query, acc)
}

func UpdateAccount(db *sqlx.DB, acc models.Account) (models.Account, error) {
	acc.UpdatedAt = time.Now()

	query := `UPDATE accounts SET account_name = :account_name, account_type = :account_type, local_share = :local_share, notes = :notes, is_active = :is_active, updated_at = :updated_at 
			  WHERE id = :id`

	return exercuteQuery(db, query, acc)
}

func DeactivateAccount(db *sqlx.DB, id string) (models.Account, error) {
	acc := models.Account{
		ID: id,
		IsActive: false,
		UpdatedAt: time.Now(),
	}

	query := `UPDATE accounts SET is_active = :is_active, updated_at = :updated_at
			  WHERE id = :id`

	return exercuteQuery(db, query, acc)
}

func GetAccount(db *sqlx.DB, id string) (models.Account, error) {
	var acc models.Account
	err := db.Get(&acc, "SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return models.Account{}, err
	}

	return acc, nil
}

func GetAllAccounts(db *sqlx.DB) ([]models.Account, error) {
	var accs []models.Account
	err := db.Select(&accs, "SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	return accs, nil
}