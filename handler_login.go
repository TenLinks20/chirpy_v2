package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request)  {
	type parameters struct {
		Email string 	`json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
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

	if p.ExpiresInSeconds == nil || *p.ExpiresInSeconds > 3600 {
		defaultVal := 3600
		p.ExpiresInSeconds = &defaultVal
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(*p.ExpiresInSeconds)*time.Second) 
		if err != nil {
			respondWithErr(w, 500, "unable to create access token")
			return
		}
	
	respBody := struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}
	respondWithJSON(w, 200, &respBody)
}