package main

import (
	"encoding/json"
	"net/http"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/TenLinks20/chirpy_v2/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request)  {
	
	type params struct {
		Email string
		Password string
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithErr(w, 400, "invalid structure")
		return
	}

	hash, err := auth.HashPassword(p.Password)
	if err != nil {
		respondWithErr(w, 500, "unable to use password")
		return
	}

	dbParams := database.CreateUserParams{
		Email: p.Email,
		HashedPassword: hash,
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), dbParams) 
	if err != nil {
		respondWithErr(w, 500, "unable to create user")
		return
	}

	user := dbToAPIUser(&dbUser)
	respondWithJSON(w, 201, &user)

}