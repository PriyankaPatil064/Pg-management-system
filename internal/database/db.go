package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "pg_management"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to Database successfully!")
}

func InitSchema() {
	createRoomsTable := `
	CREATE TABLE IF NOT EXISTS rooms (
		id SERIAL PRIMARY KEY,
		room_number VARCHAR(50) NOT NULL UNIQUE,
		capacity INT NOT NULL CHECK (capacity > 0),
		occupancy INT DEFAULT 0 CHECK (occupancy >= 0),
		price DECIMAL(10, 2) NOT NULL CHECK (price > 0)
	);`

	createGuestsTable := `
	CREATE TABLE IF NOT EXISTS guests (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(20),
		room_id INT REFERENCES rooms(id),
		join_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createPaymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		guest_id INT REFERENCES guests(id),
		amount DECIMAL(10, 2) NOT NULL,
		payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		payment_method VARCHAR(50)
	);`

	if _, err := DB.Exec(createRoomsTable); err != nil {
		log.Fatal("Failed to create rooms table:", err)
	}

	// Cleanup existing empty or invalid entries
	cleanupQuery := `DELETE FROM rooms WHERE room_number = '' OR capacity <= 0 OR price <= 0;`
	if _, err := DB.Exec(cleanupQuery); err != nil {
		log.Printf("Warning: Failed to cleanup invalid rooms: %v", err)
	} else {
		log.Println("Cleaned up existing invalid room entries.")
	}

	if _, err := DB.Exec(createGuestsTable); err != nil {
		log.Fatal("Failed to create guests table:", err)
	}

	if _, err := DB.Exec(createPaymentsTable); err != nil {
		log.Fatal("Failed to create payments table:", err)
	}

	log.Println("Database schema initialized successfully!")
}
