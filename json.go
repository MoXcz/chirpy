package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respJSON(w http.ResponseWriter, code int, jsonResp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(jsonResp)
	if err != nil {
		log.Printf("Error marshalling JSON %s\n", err)
	}
	w.Write(data)
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
