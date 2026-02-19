package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	GoogleID  string    `json:"google_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
