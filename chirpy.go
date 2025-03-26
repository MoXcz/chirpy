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

func (cfg *apiConfig) handlerCreateChirpy(w http.ResponseWriter, r *http.Request) {
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

	chirpy, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{Body: cleanBody, UserID: params.UserId})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create chirpy")
		return
	}
	respJSON(w, 201, Chirp{
		ID:        chirpy.ID,
		CreatedAt: chirpy.CreatedAt,
		UpdatedAt: chirpy.UpdatedAt,
		Body:      chirpy.Body,
		UserID:    chirpy.UserID,
	})
}
