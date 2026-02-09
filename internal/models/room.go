package models

type Room struct {
	ID         int     `json:"id"`
	RoomNumber string  `json:"room_number"`
	Capacity   int     `json:"capacity"`
	Occupancy  int     `json:"occupancy"`
	Price      float64 `json:"price"`
}
