package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MoXcz/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatalln("DB_URL variable must be defined (use a .env file)")
	}
	if platform == "" {
		log.Fatalln("PLATFORM was not defined")
	}

	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Could not open database. Err:", err)
	}
	dbQueries := database.New(dbConnection)

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, db: dbQueries, platform: platform}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
