package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type Chirpie struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
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

	cleanBody := removeProfanity(params.Body)

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

func respJSON(w http.ResponseWriter, code int, jsonResp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(jsonResp)
	if err != nil {
		log.Printf("Error marshalling JSON %s\n", err)
	}
	w.Write(data)
}

func removeProfanity(profaneBody string) string {
	words := strings.Split(profaneBody, " ")
	for i, word := range words {
		word = strings.ToLower(word)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	purifiedBody := strings.Join(words, " ")
	return purifiedBody
}
