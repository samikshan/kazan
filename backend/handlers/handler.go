package handlers

import "github.com/samikshan/kazan/backend/models"

type Handler struct {
	// repositories
	userRepo      models.UserRepo
	authRepo      models.AuthRepo
	trackRepo     models.TrackRepo
	componentRepo models.ComponentRepo
}

func New(userRepo models.UserRepo, authRepo models.AuthRepo, trackRepo models.TrackRepo, compRepo models.ComponentRepo) *Handler {
	return &Handler{
		userRepo:      userRepo,
		authRepo:      authRepo,
		trackRepo:     trackRepo,
		componentRepo: compRepo,
	}
}
