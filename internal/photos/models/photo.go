package models

type Photo struct {
	ID     string `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Base64 string `json:"base64" db:"base64"`
	UserId string `json:"userId" db:"user_id"`
}
