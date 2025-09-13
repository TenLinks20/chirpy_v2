package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
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