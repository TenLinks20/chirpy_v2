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
	secret string
	polkaKey string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("polka key must be set")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("secret must be set")
	}

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
		secret: secret,
		polkaKey: polkaKey,
	}

	// static files
	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))

	// admin handlers
	mux.Handle("GET /api/healthz", http.StripPrefix("/api", http.HandlerFunc(handlerReadiness)))
	mux.Handle("POST /api/login", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerLogin)))
	mux.Handle("POST /api/revoke", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerRevoke)))
	mux.Handle("POST /api/refresh", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerRefresh)))
	mux.Handle("GET /admin/metrics", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerMetrics)))
	mux.Handle("POST /admin/reset", http.StripPrefix("/admin", http.HandlerFunc(apiCfg.handlerReset)))

	// chirp handlers
	mux.Handle("POST /api/chirps", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerNewChirp))) 
	mux.Handle("GET /api/chirps", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerGetAllChirps)))
	mux.Handle("GET /api/chirps/{chirpID}", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerGetChirp)))
	mux.Handle("DELETE /api/chirps/{chirpID}", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerDeleteChirp)))

	// user handlers
	mux.Handle("POST /api/users", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerNewUser)))
	mux.Handle("PUT /api/users", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerUpdateUser)))

	// webhooks
	mux.Handle("POST /api/polka/webhooks", http.StripPrefix("/api", http.HandlerFunc(apiCfg.handlerUpgradeUser)))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}