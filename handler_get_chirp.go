package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request)  {
	id := r.PathValue("chirpID")
	if id == "" {
		respondWithErr(w, 400, "no id given")
		return
	}
	chirpID, err := uuid.Parse(id)
	if err != nil {
		respondWithErr(w, 400, "invalid id format")
		return
	}
	

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID) 
	if err != nil {
		respondWithErr(w, 404, "unable to retreive data")
		return
	}
	chirp := dbToAPIChirp(&dbChirp)
	respondWithJSON(w, 200, &chirp)
}