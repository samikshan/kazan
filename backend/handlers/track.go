package handlers

import (
	"net/http"

	"github.com/samikshan/kazan/backend/models"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func (h *Handler) NewTrack(c echo.Context) error {
	addr, err := walletAddrFromReq(c)
	if err != nil {
		log.WithError(err).Error("failed to get user id")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Failed to validate request sender",
		}
	}

	u, err := h.userRepo.GetByWalletAddr(addr)
	if err != nil {
		log.WithError(err).Error("failed to retrieve user for wallet address")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate request sender",
		}
	}

	if u == nil {
		log.WithError(err).WithField("walletAddress", addr).Error("no user for wallet address found")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "invalid request sender",
		}
	}

	type newTrackReq struct {
		CID           string
		Title         string
		ParentTrackID uint
		Instruments   []string
	}
	req := new(newTrackReq)
	if err := c.Bind(req); err != nil {
		log.WithError(err).Error("failed to bind request body to new track request structure")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed to process request",
		}
	}

	parentTrack, err := h.trackRepo.GetByTrackID(req.ParentTrackID)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to add track information",
		}
	}

	var instruments []models.Instrument
	if parentTrack != nil {
		instruments = parentTrack.Instruments
		log.Info(parentTrack)
	}

	for _, insName := range req.Instruments {
		ins, err := h.instrumentRepo.GetByName(insName)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "failed to add track information",
			}
		} else if ins == nil {
			log.Errorln("instrument name not found: ", insName)
			i := &models.Instrument{Name: insName}
			if err = h.instrumentRepo.Create(i); err != nil {
				log.Errorln("failed to add new instrument: ", insName)
				continue
			}
			instruments = append(instruments, *i)
		} else {
			instruments = append(instruments, *ins)
		}
	}

	t := models.Track{
		CID:           req.CID,
		Instruments:   instruments,
		Title:         req.Title,
		ParentTrackID: req.ParentTrackID,
	}

	if err := h.trackRepo.Create(&t, u); err != nil {
		log.WithError(err).Error("failed to persist new track information")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to store new track information",
		}
	}

	return c.JSON(http.StatusOK, t)
}

func (h *Handler) GetUserFeed(c echo.Context) error {
	addr, err := walletAddrFromReq(c)
	if err != nil {
		log.WithError(err).Error("failed to get user id")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Failed to validate request sender",
		}
	}

	u, err := h.userRepo.GetByWalletAddr(addr)
	if err != nil {
		log.WithError(err).Error("failed to retrieve user for wallet address")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate request sender",
		}
	}

	if u == nil {
		log.WithError(err).WithField("walletAddress", addr).Error("no user for wallet address found")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "invalid request sender",
		}
	}

	instruments := make([]string, 0)
	for _, instrument := range u.Instruments {
		instruments = append(instruments, instrument.Name)
	}

	tracksFiltered, err := h.trackRepo.GetTracksByInstrument(instruments)
	if err != nil {
		log.WithError(err).Error("failed to get tracks for instruments")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get tracks",
		}
	}

	return c.JSON(http.StatusOK, tracksFiltered)
}
