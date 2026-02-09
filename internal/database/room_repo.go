package database

import (
	"database/sql"
	"pg-management-system/internal/models"
)

func CreateRoom(room *models.Room) error {
	query := `INSERT INTO rooms (room_number, capacity, occupancy, price) 
			  VALUES ($1, $2, $3, $4) RETURNING id`
	
	return DB.QueryRow(query, room.RoomNumber, room.Capacity, room.Occupancy, room.Price).Scan(&room.ID)
}

func GetRoomByID(id int) (*models.Room, error) {
	room := &models.Room{}
	query := `SELECT id, room_number, capacity, occupancy, price FROM rooms WHERE id = $1`
	
	err := DB.QueryRow(query, id).Scan(&room.ID, &room.RoomNumber, &room.Capacity, &room.Occupancy, &room.Price)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func GetAllRooms() ([]models.Room, error) {
	rows, err := DB.Query(`SELECT id, room_number, capacity, occupancy, price FROM rooms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.RoomNumber, &room.Capacity, &room.Occupancy, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func UpdateRoom(id int, room *models.Room) error {
	query := `UPDATE rooms SET room_number=$1, capacity=$2, occupancy=$3, price=$4 WHERE id=$5`
	
	res, err := DB.Exec(query, room.RoomNumber, room.Capacity, room.Occupancy, room.Price, id)
	if err != nil {
		return err
	}
	
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteRoom(id int) error {
	query := `DELETE FROM rooms WHERE id=$1`
	
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
