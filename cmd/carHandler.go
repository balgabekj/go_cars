package main

import (
	"database/sql"
	"encoding/json"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/balgabekj/go_car/pkg/validator"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
func (app *application) getAllCarHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Brand   string
		MinYear int
		MaxYear int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Brand = app.readStrings(qs, "brand", "")
	input.MinYear = app.readInt(qs, "minyear", 0, v)
	input.MaxYear = app.readInt(qs, "maxyear", 0, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "brand", "year", "-id", "-brand", "-year"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	cars, metadata, err := app.models.Cars.GetAll(input.Brand, input.MinYear, input.MaxYear, input.Filters)

	err = app.writeJSON(w, http.StatusOK, envelope{"cars": cars, "metadata": metadata}, nil)

	if err != nil {
		app.models.Cars.ErrorLog.Println(err)
		http.Error(w, "Error retrieving cars", http.StatusInternalServerError)
		return
	}
	app.respondWithJson(w, http.StatusOK, cars)
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
	car.ID, _ = strconv.Atoi(id)

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

//func (app *application) deleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
//	// Извлечь ID элемента меню из URL
//	vars := mux.Vars(r)
//	param := vars["menuId"]
//	id, err := strconv.Atoi(param)
//	if err != nil || id < 1 {
//		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
//		return
//	}
//
//	// Проверить разрешение DELETE для пользователя
//	user := app.contextGetUser(r)
//
//
//	// Выполнить удаление элемента меню
//	err = app.models.Menu.Delete(id)
//	if err != nil {
//		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete menu item")
//		return
//	}
//
//	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
//}

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJson(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
