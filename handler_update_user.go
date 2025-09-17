package main

import (
	"encoding/json"
	"net/http"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/TenLinks20/chirpy_v2/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request)  {
	authtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, 401, "no access token")
		return
	}

	userID, err := auth.ValidateJWT(authtoken, cfg.secret)
	if err != nil {
		respondWithJSON(w, 401, "no access permitted")
		return
	}

	type parameters struct {
		Email string
		Password string
	}

	var p parameters
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithErr(w, 500, "unable to update user")
		return
	}

	hash, err := auth.HashPassword(p.Password)
	if err != nil {
		respondWithErr(w, 500, "unable to update user")
		return
	}

	dbParams := database.UpdateUserParams{
		ID: userID,
		Email: p.Email,
		HashedPassword: hash,
	}
	dbUser, err := cfg.dbQueries.UpdateUser(r.Context(), dbParams)
	if err != nil {
		respondWithErr(w, 500, "unable to update user"+err.Error())
		return
	}

	user := dbToAPIUser(&dbUser)
	respondWithJSON(w, 200, &user)
}