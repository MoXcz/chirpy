package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MoXcz/chirpy/internal/auth"
	"github.com/MoXcz/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatdeAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func newUser(id uuid.UUID, created_at, updated_at time.Time, email string) User {
	return User{
		ID:        id,
		CreatdeAt: created_at,
		UpdatedAt: updated_at,
		Email:     email,
	}
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respError(w, http.StatusBadRequest, err, "Error: Could not hash password")
		return
	}

	retUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create user")
	}

	user := newUser(retUser.ID, retUser.CreatedAt, retUser.UpdatedAt, retUser.Email)
	respJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	retUser, err := cfg.db.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(retUser.HashedPassword, params.Password)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Incorrect email or password")
		return
	}

	expirationTime := 3600 // 1 hour
	if params.ExpiresInSeconds != 0 {
		expirationTime = params.ExpiresInSeconds
	}

	user := newUser(retUser.ID, retUser.CreatedAt, retUser.UpdatedAt, retUser.Email)
	user.Token, err = auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(expirationTime))
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Could not create token")
		return
	}
	respJSON(w, http.StatusOK, user)
}
