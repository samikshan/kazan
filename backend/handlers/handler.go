package handlers

import "github.com/samikshan/kazan/backend/models"

type Handler struct {
	// repositories
	userRepo models.UserRepo
	authRepo models.AuthRepo
}

func New(userRepo models.UserRepo, authRepo models.AuthRepo) *Handler {
	return &Handler{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
