package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) GetProfile(c echo.Context) error {
	type Profile struct {
		DisplayName    string   `json:"displayName"`
		Username       string   `json:"username"`
		Instruments    []string `json:"instruments"`
		FollowingCount uint     `json:"followingCount"`
		FollowerCount  uint     `json:"followerCount"`
		JamCount       uint     `json:"jamsCount"`
		TrackCount     uint     `json:"tracksCount"`
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Error("userid should be an integer")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user id",
		}
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).Errorln("failed to get profile with ID: ", id)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get profile",
		}
	}

	if user == nil {
		log.WithError(err).WithField("id", id).Error("failed to get profile")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "invalid profile id",
		}
	}

	p := &Profile{
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}

	return c.JSON(http.StatusOK, p)
}
