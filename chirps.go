package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/MoXcz/chirpy/internal/auth"
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

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	cleanBody, err := validateChirp(params.Body)
	if err != nil {
		respError(w, http.StatusBadRequest, err, err.Error())
	}

	respChirp, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{Body: cleanBody, UserID: userId})
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not create chirpy")
		return
	}

	chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)

	respJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	respChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not get chirps")
	}

	authorID := uuid.Nil
	rawAuthorID := r.URL.Query().Get("author_id")
	if rawAuthorID != "" {
		authorID, err = uuid.Parse(rawAuthorID)
		if err != nil {
			respError(w, http.StatusBadRequest, err, "Error: Invalid author id")
			return
		}
	}

	sortOrder := r.URL.Query().Get("sort")

	chirps := []Chirp{}
	for _, respChirp := range respChirps {
		if authorID != uuid.Nil && respChirp.UserID != authorID {
			continue
		}
		chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)
		chirps = append(chirps, chirp)
	}

	// query returns chirps with "asc" order
	if sortOrder == "desc" {
		// chirps[i] has to be "newer" than chirps[j] for it to be true, meaning
		// that "i" comes before "j" -> descending order
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Compare(chirps[j].CreatedAt) == 1
		})
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
		respError(w, http.StatusNotFound, err, "Error: Could not get chirp")
		return
	}

	chirp := newChirp(respChirp.ID, respChirp.CreatedAt, respChirp.UpdatedAt, respChirp.Body, respChirp.UserID)
	respJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respError(w, http.StatusBadRequest, err, "Error: Could not parse ID")
		return
	}

	respChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respError(w, http.StatusNotFound, err, "Error: Could not get chirp")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Invalid Authorization header")
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respError(w, http.StatusUnauthorized, err, "Invalid JWT token")
		return
	}

	if respChirp.UserID != userId {
		respError(w, http.StatusForbidden, err, "Error: Not valid user, you cannot delete this chirp")
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), respChirp.ID)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
