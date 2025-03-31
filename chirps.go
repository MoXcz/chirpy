package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MoXcz/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func newChirp(id uuid.UUID, createdAt, updatedAt time.Time, body string, userId uuid.UUID) Chirp {
	return Chirp{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Body:      body,
		UserID:    userId,
	}
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	cleanBody, err := validateChirp(params.Body)
	if err != nil {
		respError(w, http.StatusBadRequest, err, err.Error())
	}

	respChirp, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{Body: cleanBody, UserID: params.UserId})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create chirpy")
		return
	}

	chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)

	respJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	respChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not get chirps")
	}
	chirps := []Chirp{}
	for _, respChirp := range respChirps {
		chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)
		chirps = append(chirps, chirp)
	}
	respJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respError(w, http.StatusBadRequest, err, "Error: Could not parse ID")
		return
	}

	respChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not get chirp")
	}

	chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)
	respJSON(w, http.StatusOK, chirp)
}
