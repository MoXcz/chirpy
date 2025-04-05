package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/MoXcz/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
	apiKey         string
}

func newConfig() apiConfig {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	apiKey := os.Getenv("POLKA_KEY")
	if dbURL == "" {
		log.Fatalln("DB_URL variable must be defined (use a .env file)")
	}
	if platform == "" {
		log.Fatalln("PLATFORM was not defined")
	}
	if tokenSecret == "" {
		log.Fatalln("TOKEN_SECRET was not defined")
	}
	if apiKey == "" {
		log.Fatalln("POLAK_KEY was not defined")
	}

	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Could not open database. Err:", err)
	}
	dbQueries := database.New(dbConnection)
	return apiConfig{fileserverHits: atomic.Int32{}, db: dbQueries, platform: platform, tokenSecret: tokenSecret, apiKey: apiKey}
}
