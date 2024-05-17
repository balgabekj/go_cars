package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car data from request body
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		app.models.Categories.ErrorLog.Println(err)
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	// Insert car into the database
	err = app.models.Categories.InsertCategory(&category)
	if err != nil {
		app.models.Categories.ErrorLog.Println(err)
		http.Error(w, "Error creating category", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (app *application) getCarByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from URL parameters
	params := mux.Vars(r)
	categoryName := params["categoryName"]

	// Retrieve car from the database
	car, err := app.models.Cars.GetByCategory(categoryName)
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

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем имя категории из параметров URL
	params := mux.Vars(r)
	categoryName := params["categoryName"]

	// Извлекаем данные категории из тела запроса
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		app.models.Categories.ErrorLog.Println(err)
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	// Обновляем категорию в базе данных
	err = app.models.Categories.UpdateCategory(categoryName, &category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		app.models.Categories.ErrorLog.Println(err)
		http.Error(w, "Error updating category", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем имя категории из параметров URL
	params := mux.Vars(r)
	categoryName := params["categoryName"]
	// Удаляем категорию из базы данных
	err := app.models.Categories.DeleteCategory(categoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		app.models.Categories.ErrorLog.Println(err)
		http.Error(w, "Error deleting category", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}
