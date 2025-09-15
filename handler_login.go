package main

import (
	"encoding/json"
	"net/http"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request)  {
	type parameters struct {
		Email string 	`json:"email"`
		Password string `json:"password"`
	}

	var p parameters
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithErr(w, 500, "invalid input")
		return
	}

	dbUser, err := cfg.dbQueries.GetUserWithEmail(r.Context(), p.Email)
	if err != nil {
		respondWithErr(w, 401, "unable to find user")
		return
	}

	if err := auth.CheckHashPassword(dbUser.HashedPassword, p.Password); err != nil {
		respondWithErr(w, 401, "incorrect email or password")
		return
	}

	user := dbToAPIUser(&dbUser)
	respondWithJSON(w, 200, &user)
}