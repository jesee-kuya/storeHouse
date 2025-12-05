package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeReceiptQuery(db *sqlx.DB, query string, receipt models.Receipt) (models.Receipt, error) {
	_, err := db.NamedExec(query, receipt)
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

func CreateReceipt(db *sqlx.DB, receipt models.Receipt) (models.Receipt, error) {
	receipt.ID = uuid.New().String()
	receipt.CreatedAt = time.Now()
	receipt.UpdatedAt = time.Now()

	query := `INSERT INTO receipts (id, transaction_id, income_account, amount, created_at, updated_at)
              VALUES (:id, :transaction_id, :income_account_id, :amount, :created_at, :updated_at)`

	return executeReceiptQuery(db, query, receipt)
}

func UpdateReceipt(db *sqlx.DB, receipt models.Receipt) (models.Receipt, error) {
	receipt.UpdatedAt = time.Now()

	query := `UPDATE receipts SET income_account = :income_account_id, amount = :amount, updated_at = :updated_at 
			  WHERE id = :id`

	return executeReceiptQuery(db, query, receipt)
}

func DeleteReceipt(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM receipts WHERE id = $1", id)
	return err
}

func GetReceipt(db *sqlx.DB, id string) (models.Receipt, error) {
	var receipt models.Receipt
	err := db.Get(&receipt, "SELECT * FROM receipts WHERE id = $1", id)
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

func GetReceiptByTransaction(db *sqlx.DB, transactionID string) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := db.Select(&receipts, "SELECT * FROM receipts WHERE transaction_id = $1 ORDER BY created_at DESC", transactionID)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func GetReceiptByAccount(db *sqlx.DB, accountID string) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := db.Select(&receipts, "SELECT * FROM receipts WHERE income_account = $1 ORDER BY created_at DESC", accountID)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func GetReceiptByDateRange(db *sqlx.DB, startDate, endDate time.Time) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := db.Select(&receipts, "SELECT * FROM receipts WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func GetAllReceipts(db *sqlx.DB) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := db.Select(&receipts, "SELECT * FROM receipts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func GetTotalReceiptsByAccount(db *sqlx.DB, accountID string) (float64, error) {
	var total float64
	err := db.Get(&total, "SELECT COALESCE(SUM(amount), 0) FROM receipts WHERE income_account = $1", accountID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetTotalReceiptsByDateRange(db *sqlx.DB, accountID string, startDate, endDate time.Time) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(amount), 0) FROM receipts WHERE income_account = $1 AND created_at BETWEEN $2 AND $3"
	err := db.Get(&total, query, accountID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	return total, nil
}
