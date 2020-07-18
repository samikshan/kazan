package handlers

import (
	"net/http"

	"github.com/samikshan/kazan/backend/models"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func (h *Handler) SetNewAuth(c echo.Context) error {
	auth := new(models.Auth)
	if err := c.Bind(auth); err != nil {
		log.WithError(err).Error("failed to bind request body to new auth request")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed to process auth request",
		}
	}

	if err := h.authRepo.Create(auth); err != nil {
		log.WithError(err).Error("signing up new user failed")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "signing up new user failed",
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetAuth(c echo.Context) error {
	type getAuthReq struct {
		LookupKey string
	}

	req := new(getAuthReq)
	if err := c.Bind(req); err != nil {
		log.WithError(err).Error("failed to bind request body")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing field: lookupKey",
		}
	}

	auth, err := h.authRepo.GetByLookupKey(req.LookupKey)
	if err != nil {
		log.WithError(err).Errorln("auth lookup failed with key: ", req.LookupKey)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to lookup auth",
		}
	}

	if auth == nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
		}
	}

	return c.JSON(http.StatusOK, auth)
}
