package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeExpenditureQuery(db *sqlx.DB, query string, exp models.Expenditure) (models.Expenditure, error) {
	_, err := db.NamedExec(query, exp)
	if err != nil {
		return models.Expenditure{}, err
	}

	return exp, nil
}

func CreateExpenditure(db *sqlx.DB, exp models.Expenditure) (models.Expenditure, error) {
	exp.ID = uuid.New().String()
	exp.CreatedAt = time.Now()
	exp.UpdatedAt = time.Now()

	query := `INSERT INTO expenditures (id, transaction_id, perticulars, bank_account, amount, created_at, updated_at)
              VALUES (:id, :transaction_id, :particulars, :bank_account_id, :amount, :created_at, :updated_at)`

	return executeExpenditureQuery(db, query, exp)
}

func UpdateExpenditure(db *sqlx.DB, exp models.Expenditure) (models.Expenditure, error) {
	exp.UpdatedAt = time.Now()

	query := `UPDATE expenditures SET perticulars = :particulars, bank_account = :bank_account_id, amount = :amount, updated_at = :updated_at 
			  WHERE id = :id`

	return executeExpenditureQuery(db, query, exp)
}

func DeleteExpenditure(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM expenditures WHERE id = $1", id)
	return err
}

func GetExpenditure(db *sqlx.DB, id string) (models.Expenditure, error) {
	var exp models.Expenditure
	err := db.Get(&exp, "SELECT * FROM expenditures WHERE id = $1", id)
	if err != nil {
		return models.Expenditure{}, err
	}

	return exp, nil
}

func GetExpenditureByTransaction(db *sqlx.DB, transactionID string) ([]models.Expenditure, error) {
	var expenses []models.Expenditure
	err := db.Select(&expenses, "SELECT * FROM expenditures WHERE transaction_id = $1 ORDER BY created_at DESC", transactionID)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func GetExpenditureByAccount(db *sqlx.DB, accountID string) ([]models.Expenditure, error) {
	var expenses []models.Expenditure
	err := db.Select(&expenses, "SELECT * FROM expenditures WHERE bank_account = $1 ORDER BY created_at DESC", accountID)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func GetExpenditureByDateRange(db *sqlx.DB, startDate, endDate time.Time) ([]models.Expenditure, error) {
	var expenses []models.Expenditure
	err := db.Select(&expenses, "SELECT * FROM expenditures WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func GetAllExpenditures(db *sqlx.DB) ([]models.Expenditure, error) {
	var expenses []models.Expenditure
	err := db.Select(&expenses, "SELECT * FROM expenditures ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
