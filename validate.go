package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type Chirpie struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, err, "Error: Could not decode parameters")
		return
	}

	if len(params.Body) > 140 {
		respError(w, http.StatusBadRequest, err, "Error: Chirpie is too long")
		return
	}

	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanBody := removeProfanity(params.Body, profanity)

	respBody := Chirpie{
		CleanedBody: cleanBody,
	}
	respJSON(w, 200, respBody)
}

func respError(w http.ResponseWriter, code int, err error, msg string) {
	if err != nil {
		log.Println(err)
	}

	type errResp struct {
		Error string `json:"error"`
	}
	errBody := errResp{
		Error: msg,
	}
	respJSON(w, code, errBody)
}

func removeProfanity(profaneBody string, profanity map[string]struct{}) string {
	words := strings.Split(profaneBody, " ")
	for i, word := range words {
		word = strings.ToLower(word)
		if _, ok := profanity[word]; ok {
			words[i] = "****"
		}
	}
	purifiedBody := strings.Join(words, " ")
	return purifiedBody
}
