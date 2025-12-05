package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeTransactionQuery(db *sqlx.DB, query string, txn models.Transaction) (models.Transaction, error) {
	_, err := db.NamedExec(query, txn)
	if err != nil {
		return models.Transaction{}, err
	}

	return txn, nil
}

func CreateTransaction(db *sqlx.DB, txn models.Transaction) (models.Transaction, error) {
	txn.ID = uuid.New().String()

	// Set default transaction date if not provided
	if txn.TransactionDate.IsZero() {
		txn.TransactionDate = time.Now()
	}

	txn.CreatedAt = time.Now()
	txn.UpdatedAt = time.Now()

	query := `INSERT INTO transactions (id, transaction_ref, transaction_date, transaction_type, amount, notes, debit_account, member, created_by, created_at, updated_at)
              VALUES (:id, :transaction_ref, :transaction_date, :transaction_type, :amount, :notes, :debit_account_id, :member_id, :created_by, :created_at, :updated_at)`

	return executeTransactionQuery(db, query, txn)
}

func UpdateTransaction(db *sqlx.DB, txn models.Transaction) (models.Transaction, error) {
	txn.UpdatedAt = time.Now()

	query := `UPDATE transactions SET transaction_ref = :transaction_ref, transaction_date = :transaction_date, transaction_type = :transaction_type, amount = :amount, notes = :notes, debit_account = :debit_account_id, member = :member_id, updated_at = :updated_at 
			  WHERE id = :id`

	return executeTransactionQuery(db, query, txn)
}

func DeleteTransaction(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM transactions WHERE id = $1", id)
	return err
}

func GetTransaction(db *sqlx.DB, id string) (models.Transaction, error) {
	var txn models.Transaction
	err := db.Get(&txn, "SELECT * FROM transactions WHERE id = $1", id)
	if err != nil {
		return models.Transaction{}, err
	}

	return txn, nil
}

func GetTransactionByRef(db *sqlx.DB, ref string) (models.Transaction, error) {
	var txn models.Transaction
	err := db.Get(&txn, "SELECT * FROM transactions WHERE transaction_ref = $1", ref)
	if err != nil {
		return models.Transaction{}, err
	}

	return txn, nil
}

func GetTransactionsByAccount(db *sqlx.DB, accountID string) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := db.Select(&txns, "SELECT * FROM transactions WHERE debit_account = $1 ORDER BY transaction_date DESC", accountID)
	if err != nil {
		return nil, err
	}

	return txns, nil
}

func GetTransactionsByMember(db *sqlx.DB, memberID string) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := db.Select(&txns, "SELECT * FROM transactions WHERE member = $1 ORDER BY transaction_date DESC", memberID)
	if err != nil {
		return nil, err
	}

	return txns, nil
}

func GetTransactionsByType(db *sqlx.DB, txnType string) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := db.Select(&txns, "SELECT * FROM transactions WHERE transaction_type = $1 ORDER BY transaction_date DESC", txnType)
	if err != nil {
		return nil, err
	}

	return txns, nil
}

func GetAllTransactions(db *sqlx.DB) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := db.Select(&txns, "SELECT * FROM transactions ORDER BY transaction_date DESC")
	if err != nil {
		return nil, err
	}

	return txns, nil
}

func GetTransactionsByDateRange(db *sqlx.DB, startDate, endDate time.Time) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := db.Select(&txns, "SELECT * FROM transactions WHERE transaction_date BETWEEN $1 AND $2 ORDER BY transaction_date DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}

	return txns, nil
}
