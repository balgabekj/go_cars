package main

import (
	"database/sql"
	"encoding/json"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createCarHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car data from request body
	var car model.Car
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	// Insert car into the database
	err = app.models.Cars.Insert(&car)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error creating car", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

func (app *application) getCarHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Retrieve car from the database
	car, err := app.models.Cars.Get(id)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		if err == sql.ErrNoRows {
			http.Error(w, "Car not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving car", http.StatusInternalServerError)
		}
		return
	}

	// Return car as JSON response
	json.NewEncoder(w).Encode(car)
}

func (app *application) updateCarHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Extract car data from request body
	var car model.Car
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	// Set car ID for update
	car.ID = id

	// Update car in the database
	err = app.models.Cars.Update(&car)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error updating car", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(car)
}

func (app *application) deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Delete car from the database
	err := app.models.Cars.Delete(id)
	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error deleting car", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}
