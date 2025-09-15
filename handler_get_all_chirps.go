package main

import "net/http"

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request)  {
	chirps := []APIChirp{}

	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithErr(w, 500, "unable to retreive data")
		return
	}

	for _, dbChirp := range dbChirps {
		chirp := dbToAPIChirp(&dbChirp)
		chirps = append(chirps, chirp)
	}

	if len(chirps) != len(dbChirps) {
		respondWithErr(w, 500, "bad data retreived")
		return
	}

	respondWithJSON(w, 200, &chirps)
}