package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func application(usersModel UsersModel) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		u, err := usersModel.GetUsers()

		if err != nil {
			fmt.Println("Error occurred:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		u, err := usersModel.GetUser(id)
		if err != nil {
			switch err.(type) {
			case UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Not Found"))
				return
			}
			fmt.Println("Error occurred:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := usersModel.DeleteUser(id)
		if err != nil {
			fmt.Println("Error occurred:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Println("Error occurred:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		newUser, err := usersModel.CreateUser(user)
		if err != nil {
			switch err.(type) {
			case UniquenessViolationError:
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Conflict"))
				return
			}
			fmt.Println("Error occurred:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newUser)
	})

	return r
}
