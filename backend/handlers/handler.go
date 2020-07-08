package handlers

import "github.com/samikshan/kazan/backend/models"

type Handler struct {
	// repositories
	userRepo models.UserRepo
}

func New(userRepo models.UserRepo) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}
