package database

import (
	"database/sql"
	"errors"
	"shop_backend/models"
)

func DBInitUsers() {
	gDB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		passwordHash TEXT,
		role TEXT
	)
	`)
}

var ErrUserNotFound = errors.New("user not found")

func CreateUser(username string, passwordHash []byte, role string) error {
	_, err := gDB.Exec("INSERT INTO users (username, passwordHash, role) VALUES (?, ?, ?)", username, passwordHash, role)
	return err
}

func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	err := gDB.QueryRow("SELECT id, username, role, passwordHash FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Role, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UpdatePasswordHash(userId int, passwordHash []byte) error {
	_, err := gDB.Exec("UPDATE passwordHash FROM users WHERE id = ?", passwordHash, userId)
	return err
}
