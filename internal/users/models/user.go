package models

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"title" db:"title"`
	Password string `json:"base64" db:"base64"`
}
