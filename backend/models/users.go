package models

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	PasswordHash	string
}
