package user

import "time"

type User struct {
	ID           string    `json:"id"`
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
}
