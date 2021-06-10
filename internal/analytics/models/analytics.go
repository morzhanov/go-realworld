package models

type AnalyticsEntry struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}
