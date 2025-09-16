package main

import (
	"net/http"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request)  {
	authValue, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, 401, "no valid token provided")
		return
	}

	if err := cfg.dbQueries.RevokeRefreshToken(r.Context(), authValue); err != nil {
		respondWithErr(w, 500, "unable to revoke token")
		return
	}

	w.WriteHeader(204)
}