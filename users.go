package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/MoXcz/chirpy/internal/auth"
	"github.com/MoXcz/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatdeAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func newUser(id uuid.UUID, created_at, updated_at time.Time, email string, isChirpyRed bool) User {
	return User{
		ID:          id,
		CreatdeAt:   created_at,
		UpdatedAt:   updated_at,
		Email:       email,
		IsChirpyRed: isChirpyRed,
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

	user := newUser(retUser.ID, retUser.CreatedAt, retUser.UpdatedAt, retUser.Email, retUser.IsChirpyRed)
	respJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respError(w, 401, err, "Invalid Authorization header")
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respError(w, 401, err, "Invalid JWT token")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respError(w, http.StatusBadRequest, err, "Error: Could not hash password")
		return
	}
	retUser, err := cfg.db.UpdateUserFromId(r.Context(), database.UpdateUserFromIdParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not update user")
	}

	user := newUser(retUser.ID, retUser.CreatedAt, retUser.UpdatedAt, retUser.Email, retUser.IsChirpyRed)
	respJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	expirationTime := time.Hour

	user := newUser(retUser.ID, retUser.CreatedAt, retUser.UpdatedAt, retUser.Email, retUser.IsChirpyRed)
	user.Token, err = auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(expirationTime))
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Could not create JWT")
		return
	}
	user.RefreshToken, err = auth.MakeRefreshToken()
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Could not create token")
		return
	}
	_, err = cfg.db.CreateToken(r.Context(), database.CreateTokenParams{
		Token:     user.RefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Could not store token")
		return
	}
	respJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respError(w, 401, err, "Invalid Authorization header")
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respError(w, 401, err, "Invalid token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Could not create JWT")
		return
	}

	respJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respError(w, http.StatusBadRequest, err, "Invalid Authorization header")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Invalid token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respError(w, 401, err, "Missing API key in Authorization header")
		return
	}

	if apiKey != cfg.apiKey {
		respError(w, 401, err, "Missing API key in Authorization header")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respError(w, http.StatusNotFound, err, "Error: Could not find user")
			return
		}
		respError(w, http.StatusNotFound, err, "Error: Could not upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
