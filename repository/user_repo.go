package repository

import (
	"storeHouse/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func executeUserQuery(db *sqlx.DB, query string, user models.User) (models.User, error) {
	_, err := db.NamedExec(query, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateUser(db *sqlx.DB, user models.User) (models.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `INSERT INTO users (id, username, email, password_hash, full_name, role, phone_number, is_active, created_at, updated_at)
              VALUES (:id, :username, :email, :password_hash, :full_name, :role, :phone_number, :is_active, :created_at, :updated_at)`

	return executeUserQuery(db, query, user)
}

func UpdateUser(db *sqlx.DB, user models.User) (models.User, error) {
	user.UpdatedAt = time.Now()

	query := `UPDATE users SET username = :username, email = :email, full_name = :full_name, role = :role, phone_number = :phone_number, is_active = :is_active, updated_at = :updated_at 
			  WHERE id = :id`

	return executeUserQuery(db, query, user)
}

func DeactivateUser(db *sqlx.DB, id string) (models.User, error) {
	user := models.User{
		ID:        id,
		IsActive:  false,
		UpdatedAt: time.Now(),
	}

	query := `UPDATE users SET is_active = :is_active, updated_at = :updated_at WHERE id = :id`

	return executeUserQuery(db, query, user)
}

func GetUser(db *sqlx.DB, id string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByUsername(db *sqlx.DB, username string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *sqlx.DB, email string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UpdateLastLogin(db *sqlx.DB, id string) error {
	_, err := db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", time.Now(), id)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersByRole(db *sqlx.DB, role string) ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM users WHERE role = $1 ORDER BY created_at DESC", role)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetActiveUsers(db *sqlx.DB) ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM users WHERE is_active = true ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	return users, nil
}
