package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	fmt.Println("Running")
	r := mux.NewRouter()
	// Cars
	cars := r.PathPrefix("/api/v1").Subrouter()

	cars.HandleFunc("/cars", app.createCarHandler).Methods("POST")
	cars.HandleFunc("/cars/{id}", app.getCarHandler).Methods("GET")
	cars.HandleFunc("/cars", app.getAllCarHandler).Methods("GET")
	cars.HandleFunc("/cars/{id}", app.updateCarHandler).Methods("PUT")
	cars.HandleFunc("/cars/{id}", app.requirePermissions("cars:write", app.deleteCarHandler)).Methods("DELETE")

	// Category
	cars.HandleFunc("/category/{categoryName}/cars", app.getCarByCategoryHandler).Methods("GET")
	cars.HandleFunc("/category", app.createCategoryHandler).Methods("POST")
	cars.HandleFunc("/category/{name}", app.updateCategoryHandler).Methods("PUT")
	cars.HandleFunc("/category/{name}", app.deleteCategoryHandler).Methods("DELETE")
	//Users
	users := r.PathPrefix("/api/v1").Subrouter()

	users.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)

}
