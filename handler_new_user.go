package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request)  {
	
	type params struct {
		Email string
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithErr(w, 400, "invalid structure")
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), p.Email) 
	if err != nil {
		respondWithErr(w, 500, "unable to create user")
		return
	}

	user := dbToAPIUser(&dbUser)
	respondWithJSON(w, 201, &user)


}