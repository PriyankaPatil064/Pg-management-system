package database

import (
	"database/sql"
	"pg-management-system/internal/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, name, google_id, role, created_at FROM users WHERE email = $1`
	row := DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.GoogleID, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	if user.Role == "" {
		user.Role = "user"
	}
	query := `INSERT INTO users (email, name, google_id, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := DB.QueryRow(query, user.Email, user.Name, user.GoogleID, user.Role).Scan(&user.ID, &user.CreatedAt)
	return err
}

func UpdateUserByGoogleID(googleID string, name string) error {
	query := `UPDATE users SET name = $1 WHERE google_id = $2`
	_, err := DB.Exec(query, name, googleID)
	return err
}
