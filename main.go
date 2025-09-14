package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/TenLinks20/chirpy_v2/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("db conn string needed")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("unable to open db: %s", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: dbQueries,
	}

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))

	mux.Handle("GET /api/healthz", http.StripPrefix("/api", http.HandlerFunc(handlerReadiness)))
	mux.Handle("POST /api/validate_chirp", http.StripPrefix("/api", http.HandlerFunc(handlerValidateChirps)))
	mux.Handle("GET /admin/metrics", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerMetrics)))
	mux.Handle("POST /admin/reset", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerReset)))


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}