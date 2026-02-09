package database

import (
	"time"

	"pg-management-system/internal/models"
)

func CreatePayment(payment *models.Payment) error {
	query := `
		INSERT INTO payments (guest_id, amount, payment_date, payment_method)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	return DB.QueryRow(
		query,
		payment.GuestID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
	).Scan(&payment.ID)
}

func GetPaymentsByGuestID(guestID int) ([]models.Payment, error) {
	query := `
		SELECT id, guest_id, amount, payment_date, payment_method
		FROM payments
		WHERE guest_id = $1
	`

	rows, err := DB.Query(query, guestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.GuestID,
			&p.Amount,
			&p.PaymentDate,
			&p.PaymentMethod,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func GetAllPayments() ([]models.Payment, error) {
	query := `
		SELECT id, guest_id, amount, payment_date, payment_method
		FROM payments
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.ID,
			&p.GuestID,
			&p.Amount,
			&p.PaymentDate,
			&p.PaymentMethod,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func GetPaymentByID(id int) (*models.Payment, error) {
	query := `
		SELECT id, guest_id, amount, payment_date, payment_method
		FROM payments
		WHERE id = $1
	`

	var p models.Payment
	err := DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.GuestID,
		&p.Amount,
		&p.PaymentDate,
		&p.PaymentMethod,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func UpdatePayment(payment *models.Payment) error {
	query := `
		UPDATE payments
		SET guest_id = $1, amount = $2, payment_date = $3, payment_method = $4
		WHERE id = $5
	`

	_, err := DB.Exec(
		query,
		payment.GuestID,
		payment.Amount,
		payment.PaymentDate,
		payment.PaymentMethod,
		payment.ID,
	)
	return err
}

func DeletePayment(id int) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}
