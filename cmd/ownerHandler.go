package main

import (
	"database/sql"
	"encoding/json"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createOwnerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract owner data from request body
	var owner model.Owner
	err := app.readJSON(w, r, &owner)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Error decoding data")
		return
	}

	// Insert owner into the database
	err = app.models.Owners.Insert(&owner)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error creating owner")
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(owner)
}

func (app *application) getOwnerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract owner ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Retrieve owner from the database
	owner, err := app.models.Owners.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.respondWithError(w, http.StatusNotFound, "Owner not found")
		} else {
			app.respondWithError(w, http.StatusInternalServerError, "Error retrieving owner")
		}
		return
	}

	// Return owner as JSON response
	app.respondWithJson(w, http.StatusOK, owner)
}

func (app *application) getAllOwnersHandler(w http.ResponseWriter, r *http.Request) {
	owners, err := app.models.Owners.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error retrieving owners")
		return
	}

	// Return owners as JSON response
	app.respondWithJson(w, http.StatusOK, owners)
}

func (app *application) updateOwnerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract owner ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Extract owner data from request body
	var owner model.Owner
	err := app.readJSON(w, r, &owner)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Error decoding data")
		return
	}

	// Set owner ID for update
	owner.ID = id

	// Update owner in the database
	err = app.models.Owners.Update(&owner)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error updating owner")
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owner)
}

func (app *application) deleteOwnerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract owner ID from URL parameters
	params := mux.Vars(r)
	id := params["id"]

	// Delete owner from the database
	err := app.models.Owners.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error deleting owner")
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}
