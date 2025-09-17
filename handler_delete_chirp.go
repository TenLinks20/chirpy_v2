package main

import (
	"net/http"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request)  {
	authtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, 401, "no token provided")
		return
	}

	userID, err := auth.ValidateJWT(authtoken, cfg.secret)
	if err != nil {
		respondWithErr(w, 401, "token not valid")
		return
	}

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithErr(w, 400, "no id in url")
		return
	}
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithErr(w, 404, "invalid id format")
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithErr(w, 404, "unable to find chirp")
		return
	}

	if userID != dbChirp.UserID {
		respondWithErr(w, 403, "no access permitted")
		return
	}
	if err := cfg.dbQueries.DeleteChirp(r.Context(), chirpUUID); err != nil {
		respondWithErr(w, 500, "unable to delete chirp")
		return
	}
	w.WriteHeader(204)
}