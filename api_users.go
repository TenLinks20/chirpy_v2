package main

import (
	"time"

	"github.com/TenLinks20/chirpy_v2/internal/database"
	"github.com/google/uuid"
)

type APIUser struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func dbToAPIUser(dbUser *database.User) APIUser {
	return APIUser{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}