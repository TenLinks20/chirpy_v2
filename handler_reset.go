package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request)  {
	if cfg.platform != "dev" {
		respondWithErr(w, 403, "no permitted access")
		return
	}

	if err := cfg.dbQueries.Reset(r.Context()); err != nil {
		respondWithErr(w, 500, "unable to reset")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
}