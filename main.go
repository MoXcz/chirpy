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
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL variable must be defined (use a .env file)")
	}

	dbConnection, err := sql.Open("postgrs", dbURL)
	if err != nil {
		log.Println("Could not open database. Err:", err)
	}
	dbQueries := database.New(dbConnection)

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, db: dbQueries}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
