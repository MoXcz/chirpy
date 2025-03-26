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

	chirp, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{Body: cleanBody, UserID: params.UserId})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create chirpy")
		return
	}
	respJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	respChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not get chirps")
	}
	chirps := []Chirp{}
	for _, chirp := range respChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respJSON(w, http.StatusOK, chirps)
}
