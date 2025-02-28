package database

import (
	"time"
	"github.com/google/uuid"
	"shop_backend/models"
)

func DBInitSessions() {
	gDB.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		sessionId TEXT PRIMARY KEY,
		userId INTEGER,
		expireTime INTEGER,
		FOREIGN KEY(userId) REFERENCES users(id)
	)
	`)
}

func CreateSession(userId int) string {
	sessionId := uuid.New().String()
	expireTime := time.Now().Add(time.Hour * 24).Unix()
	gDB.Exec("INSERT INTO sessions (sessionId, userId, expireTime) VALUES (?, ?, ?)", sessionId, userId, expireTime)
	return sessionId
}

func RemoveSession(sessionId string) error {
	_, err := gDB.Exec("DELETE FROM sessions WHERE sessionId = ?", sessionId)
	return err
}

func GetUserBySession(sessionId string) (models.User, error) {
	var session models.Session
	err := gDB.QueryRow("SELECT userId, expireTime FROM sessions WHERE sessionId = ?", sessionId).Scan(&session.UserId, &session.ExpireTime)
	if err != nil {
		return models.User{}, models.ErrInvalidSession
	}
	if session.ExpireTime < time.Now().Unix() {
		gDB.Exec("DELETE FROM sessions WHERE sessionId = ?", sessionId)
		return models.User{}, models.ErrSessionExpired
	}
	var user models.User
	err = gDB.QueryRow("SELECT id, username, role, passwordHash FROM users WHERE id = ?", session.UserId).Scan(&user.Id, &user.Username, &user.Role, &user.PasswordHash)
	if err != nil {
		return models.User{}, models.ErrInvalidUser
	}
	return user, nil
}
