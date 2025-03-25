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
