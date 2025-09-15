package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/TenLinks20/chirpy_v2/internal/database"
)

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request)  {

	authValue, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, 401, "no auth token")
		return
	}

	userID, err := auth.ValidateJWT(authValue, cfg.secret)
	if err != nil {
		respondWithErr(w, 401, "no access permitted")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithErr(w, 400, "invalid input")
		return
	}

	chirpLength := utf8.RuneCountInString(params.Body)
	if chirpLength > 140 {
		respondWithErr(w, 400, "body is too long")
		return
	}

	badWords := map[string]struct{}{}
	badWords["kerfuffle"] = struct{}{}
	badWords["sharbert"] = struct{}{}
	badWords["fornax"] = struct{}{}

	words := strings.Fields(params.Body)
	replacementStr := "****"
	for badWord := range badWords {
		for i, word := range words {
			if badWord == strings.ToLower(word) {
				words[i] = replacementStr
			}
		}
	}
	cleanedBody := strings.Join(words, " ")

	dbParams := database.CreateChirpParams{
		Body: cleanedBody,
		UserID: userID,
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithErr(w, 500, "unable to create chirp")
		return
	}
	chirp := dbToAPIChirp(&dbChirp)

	respondWithJSON(w, 201, &chirp)
}