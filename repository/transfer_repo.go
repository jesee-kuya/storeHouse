package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeTransferQuery(db *sqlx.DB, query string, transfer models.Transfer) (models.Transfer, error) {
	_, err := db.NamedExec(query, transfer)
	if err != nil {
		return models.Transfer{}, err
	}

	return transfer, nil
}

func CreateTransfer(db *sqlx.DB, transfer models.Transfer) (models.Transfer, error) {
	transfer.ID = uuid.New().String()
	transfer.CreatedAt = time.Now()
	transfer.UpdatedAt = time.Now()

	query := `INSERT INTO transfers (id, transaction_id, perticulars, credit_account, amount, created_at, updated_at)
              VALUES (:id, :transaction_id, :particulars, :credit_account_id, :amount, :created_at, :updated_at)`

	return executeTransferQuery(db, query, transfer)
}

func UpdateTransfer(db *sqlx.DB, transfer models.Transfer) (models.Transfer, error) {
	transfer.UpdatedAt = time.Now()

	query := `UPDATE transfers SET perticulars = :particulars, credit_account = :credit_account_id, amount = :amount, updated_at = :updated_at 
			  WHERE id = :id`

	return executeTransferQuery(db, query, transfer)
}

func DeleteTransfer(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM transfers WHERE id = $1", id)
	return err
}

func GetTransfer(db *sqlx.DB, id string) (models.Transfer, error) {
	var transfer models.Transfer
	err := db.Get(&transfer, "SELECT * FROM transfers WHERE id = $1", id)
	if err != nil {
		return models.Transfer{}, err
	}

	return transfer, nil
}

func GetTransferByTransaction(db *sqlx.DB, transactionID string) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := db.Select(&transfers, "SELECT * FROM transfers WHERE transaction_id = $1 ORDER BY created_at DESC", transactionID)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func GetTransferByCreditAccount(db *sqlx.DB, accountID string) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := db.Select(&transfers, "SELECT * FROM transfers WHERE credit_account = $1 ORDER BY created_at DESC", accountID)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func GetTransferByDateRange(db *sqlx.DB, startDate, endDate time.Time) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := db.Select(&transfers, "SELECT * FROM transfers WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func GetAllTransfers(db *sqlx.DB) ([]models.Transfer, error) {
	var transfers []models.Transfer
	err := db.Select(&transfers, "SELECT * FROM transfers ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func GetTotalTransfersByCreditAccount(db *sqlx.DB, accountID string) (float64, error) {
	var total float64
	err := db.Get(&total, "SELECT COALESCE(SUM(amount), 0) FROM transfers WHERE credit_account = $1", accountID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetTotalTransfersByDateRange(db *sqlx.DB, accountID string, startDate, endDate time.Time) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(amount), 0) FROM transfers WHERE credit_account = $1 AND created_at BETWEEN $2 AND $3"
	err := db.Get(&total, query, accountID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	return total, nil
}
