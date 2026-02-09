package database

import (
	"database/sql"
	"time"
	"pg-management-system/internal/models"
)

func CreateGuest(guest *models.Guest) error {
	query := `INSERT INTO guests (name, email, phone, room_id, join_date) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	if guest.JoinDate.IsZero() {
		guest.JoinDate = time.Now()
	}

	return DB.QueryRow(query, guest.Name, guest.Email, guest.Phone, guest.RoomID, guest.JoinDate).Scan(&guest.ID)
}

func GetGuestByID(id int) (*models.Guest, error) {
	guest := &models.Guest{}
	query := `SELECT id, name, email, phone, room_id, join_date FROM guests WHERE id = $1`
	
	err := DB.QueryRow(query, id).Scan(&guest.ID, &guest.Name, &guest.Email, &guest.Phone, &guest.RoomID, &guest.JoinDate)
	if err != nil {
		return nil, err
	}
	return guest, nil
}

func GetAllGuests() ([]models.Guest, error) {
	rows, err := DB.Query(`SELECT id, name, email, phone, room_id, join_date FROM guests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []models.Guest
	for rows.Next() {
		var guest models.Guest
		if err := rows.Scan(&guest.ID, &guest.Name, &guest.Email, &guest.Phone, &guest.RoomID, &guest.JoinDate); err != nil {
			return nil, err
		}
		guests = append(guests, guest)
	}
	return guests, nil
}

func UpdateGuest(id int, guest *models.Guest) error {
	query := `UPDATE guests SET name=$1, email=$2, phone=$3, room_id=$4 WHERE id=$5`
	
	res, err := DB.Exec(query, guest.Name, guest.Email, guest.Phone, guest.RoomID, id)
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

func DeleteGuest(id int) error {
	query := `DELETE FROM guests WHERE id=$1`
	
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
