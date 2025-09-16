package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
	"github.com/TenLinks20/chirpy_v2/internal/database"
	"github.com/google/uuid"
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour) 
		if err != nil {
			respondWithErr(w, 500, "unable to create access token")
			return
		}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithErr(w, 500, "unable to make refresh token")
		return
	}
	
	dbParams := database.NewRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}
	dbRefreshToken, err := cfg.dbQueries.NewRefreshToken(r.Context(), dbParams)
	if err != nil {
		respondWithErr(w, 500, "unable to store token")
		return
	}

	respBody := struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: dbRefreshToken.Token,
	}
	respondWithJSON(w, 200, &respBody)
}