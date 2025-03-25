package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatdeAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	retUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create user")
	}

	user := User{
		ID:        retUser.ID,
		CreatdeAt: retUser.CreatedAt,
		UpdatedAt: retUser.UpdatedAt,
		Email:     retUser.Email,
	}

	respJSON(w, http.StatusCreated, user)
}
