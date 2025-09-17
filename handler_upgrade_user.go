package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request)  {
	
	type polkaParams struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var p polkaParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithErr(w, 500, "unable to use given data")
		return
	}

	userID, err := uuid.Parse(p.Data.UserID)
	if err != nil {
		respondWithErr(w, 404, "id not a valid uuid")
		return
	}

	if p.Event == "user.upgraded" {
		if err := cfg.dbQueries.UpgradeUser(r.Context(), userID); err != nil {
			respondWithErr(w, 404, "unable to upgrade user")
			return
		}
	}
	w.WriteHeader(204)
}