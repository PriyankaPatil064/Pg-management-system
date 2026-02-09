package models

import "time"

type Payment struct {
	ID            int       `json:"id"`
	GuestID       int       `json:"guest_id"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	PaymentMethod string    `json:"payment_method"`
}
