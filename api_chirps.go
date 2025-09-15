package main

import (
	"time"

	"github.com/TenLinks20/chirpy_v2/internal/database"
	"github.com/google/uuid"
)

type APIChirp struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID    `json:"user_id"`
}

func dbToAPIChirp(dbChirp *database.Chirp) APIChirp {
	return APIChirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}
}