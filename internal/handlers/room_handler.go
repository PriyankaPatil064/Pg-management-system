package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pg-management-system/internal/database"
	"pg-management-system/internal/models"

	"github.com/gorilla/mux"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validation
	if room.RoomNumber == "" {
		http.Error(w, "Room number is required", http.StatusBadRequest)
		return
	}
	if room.Capacity <= 0 {
		http.Error(w, "Capacity must be greater than zero", http.StatusBadRequest)
		return
	}
	if room.Price <= 0 {
		http.Error(w, "Price must be greater than zero", http.StatusBadRequest)
		return
	}

	if err := database.CreateRoom(&room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := database.GetAllRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	room, err := database.GetRoomByID(id)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func UpdateRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.UpdateRoom(id, &room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteRoom(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
