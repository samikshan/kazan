package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/samikshan/kazan/backend/models"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) Signup(c echo.Context) error {
	// Bind request
	type signupReq struct {
		Email    string
		Password string
	}

	req := &signupReq{}
	if err := c.Bind(req); err != nil {
		log.WithError(err).Error("failed to bind request body to signup request structure")
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "failed to process signup request",
		}
	}

	log.Info(req)

	// Validate
	if len(req.Email) == 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "email cannot be empty"}
	}

	if u, err := h.userRepo.GetByEmail(req.Email); u != nil && err == nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "user already exists"}
	}

	u := &models.User{
		Email: req.Email,
	}

	if err := u.HashPassword(req.Password); err != nil {
		log.WithError(err).Error("failed to generate password hash")
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "signup failed"}
	}

	// Save user
	if err := h.userRepo.Create(u); err != nil {
		log.WithError(err).Error("failed to add user to db")
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "signup failed"}
	}

	u.Password = ""

	return c.JSON(http.StatusCreated, u)
}
