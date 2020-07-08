package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/samikshan/kazan/backend"
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

func (h *Handler) Login(c echo.Context) error {
	// Bind request
	type loginReq struct {
		Email    string
		Password string
	}

	req := &loginReq{}
	if err := c.Bind(req); err != nil {
		log.WithError(err).Error("failed to bind request body to login request structure")
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "failed to process login request",
		}
	}

	// Validate
	if len(req.Email) == 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "email cannot be empty"}
	}

	// get that awesome user who decided to log in!
	u, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		log.WithError(err).WithField("email", req.Email).Error("failed to get user from database")
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid user"}
	}

	if u == nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid email or password",
		}
	}

	if u.CheckPassword(req.Password) == false {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid email or password",
		}
	}

	if len(u.FFSToken) == 0 {
		_, ffsToken, err := u.FFSCreate()
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "failed to create ffs auth token",
			}
		}

		u.FFSToken = ffsToken
		if err := h.userRepo.Update(u); err != nil {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "failed to store ffs auth token",
			}
		}
	}

	tokens, err := generateTokenPair(u, backend.Cfg.JWT.SigningKey)
	if err != nil {
		log.WithError(err).Error("failed to generate fresh token pair")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to generate token pair",
		}
	}

	tokens["ffs_token"] = u.FFSToken

	return c.JSON(http.StatusOK, tokens)
}
