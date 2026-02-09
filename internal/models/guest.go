package models

import "time"

type Guest struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	RoomID   int       `json:"room_id"`
	JoinDate time.Time `json:"join_date"`
}
