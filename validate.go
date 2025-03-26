package main

import (
	"errors"
	"strings"
)

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

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Error: Chirpie is too long")
	}

	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := removeProfanity(body, profanity)
	return cleaned, nil
}
