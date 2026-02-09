package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pg-management-system/internal/database"
	"pg-management-system/internal/models"
	"github.com/gorilla/mux"
)

func CreateGuest(w http.ResponseWriter, r *http.Request) {
	var guest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&guest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.CreateGuest(&guest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(guest)
}

func GetAllGuests(w http.ResponseWriter, r *http.Request) {
	guests, err := database.GetAllGuests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guests)
}

func GetGuestByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	guest, err := database.GetGuestByID(id)
	if err != nil {
		http.Error(w, "Guest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guest)
}

func UpdateGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var guest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&guest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.UpdateGuest(id, &guest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guest.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guest)
}

func DeleteGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteGuest(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
