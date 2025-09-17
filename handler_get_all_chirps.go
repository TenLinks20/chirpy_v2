package main

import (
	"net/http"
	"sort"

	"github.com/TenLinks20/chirpy_v2/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request)  {

	authorID := r.URL.Query().Get("author_id")
	pattern := r.URL.Query().Get("sort")
	var dbChirps []database.Chirp
	var err error

	if authorID != "" {
		authorUUID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithErr(w, 400, "author id not valid")
			return
		}
		dbChirps, err = cfg.dbQueries.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithErr(w, 500, "unable to get chirps")
			return
		}
	} else {
			dbChirps, err = cfg.dbQueries.GetAllChirps(r.Context())
			if err != nil {
			respondWithErr(w, 500, "unable to retreive data")
			return
		}
	}
	
	chirps := []APIChirp{}
	for _, dbChirp := range dbChirps {
		chirp := dbToAPIChirp(&dbChirp)
		chirps = append(chirps, chirp)
	}

	if pattern == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, 200, &chirps)
}