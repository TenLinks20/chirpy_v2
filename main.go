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
	platform string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("platform must be set")
	}

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
		platform: platform,
	}

	// static files
	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))

	// admin handlers
	mux.Handle("GET /api/healthz", http.StripPrefix("/api", http.HandlerFunc(handlerReadiness)))
	mux.Handle("GET /admin/metrics", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerMetrics)))
	mux.Handle("POST /admin/reset", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerReset)))

	// chirp handlers
	mux.Handle("POST /api/chirps", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerNewChirp)))

	// user handlers
	mux.Handle("POST /api/users", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerNewUser)))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}