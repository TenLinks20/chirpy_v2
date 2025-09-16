package main

import (
	"net/http"
	"time"

	"github.com/TenLinks20/chirpy_v2/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request)  {
	authValue, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, 401, "no token present in header")
		return
	}

	dbRefreshToken, err := cfg.dbQueries.FindRefreshToken(r.Context(), authValue)
	if err != nil {
		respondWithErr(w, 401, "no valid token found"+err.Error())
		return
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) {
		respondWithErr(w, 401, "token expired")
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		respondWithErr(w, 401, "token revoked")
		return
	}

	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithJSON(w, 500, "unable to grant  access token")
		return
	}

	respBody := struct {Token string `json:"token"`}{Token: token}
	respondWithJSON(w, 200, &respBody)
	
}