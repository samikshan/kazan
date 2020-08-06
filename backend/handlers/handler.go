package handlers

import "github.com/samikshan/kazan/backend/models"

type Handler struct {
	// repositories
	userRepo       models.UserRepo
	authRepo       models.AuthRepo
	trackRepo      models.TrackRepo
	instrumentRepo models.InstrumentRepo
}

func New(userRepo models.UserRepo, authRepo models.AuthRepo, trackRepo models.TrackRepo, insRepo models.InstrumentRepo) *Handler {
	return &Handler{
		userRepo:       userRepo,
		authRepo:       authRepo,
		trackRepo:      trackRepo,
		instrumentRepo: insRepo,
	}
}
