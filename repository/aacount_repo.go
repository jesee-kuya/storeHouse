package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateAccount(db *sqlx.DB, acc models.Account) (models.Account, error) {
	acc.ID = uuid.New().String()
	acc.CreatedAt = time.Now()
	acc.UpdatedAt = time.Now()

	query := `INSERT INTO accounts (id, account_name, account_type, local_share, notes, is_active, created_at, updated_at)
              VALUES (:id, :account_name, :account_type, :local_share, :notes, :is_active, :created_at, :updated_at)`

	_, err := db.NamedExec(query, acc)
	if err != nil {
		return models.Account{}, err
	}

	return acc, nil
}
