package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/samikshan/kazan/backend"
	"github.com/samikshan/kazan/backend/models"
	log "github.com/sirupsen/logrus"
)

func generateTokenPair(u *models.User, signingKey string) (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.WithError(err).Error("failed to get signed token")
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = u.ID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(signingKey))
	if err != nil {
		log.WithError(err).Error("failed to get signed refresh token")
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}

func (h *Handler) Token(c echo.Context) error {
	type refreshTokenReq struct {
		RefreshToken string
	}

	refreshTokReq := refreshTokenReq{}
	if err := c.Bind(&refreshTokReq); err != nil {
		log.WithError(err).Error("failed to bind request body")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed to process request",
		}
	}

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(refreshTokReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// validate the alg with expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, echo.ErrUnauthorized
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(backend.Cfg.JWT.SigningKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := uint(claims["sub"].(float64))

		u, err := h.userRepo.GetByID(id)
		if err != nil {
			log.WithError(err).Error("failed to get user from db")
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "user not found",
			}
		}

		newTokenPair, err := generateTokenPair(u, backend.Cfg.JWT.SigningKey)
		if err != nil {
			log.WithError(err).Error("failed to generate fresh token pair")
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "failed to generate token pair",
			}
		}

		return c.JSON(http.StatusOK, newTokenPair)
	}

	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: "invalid token",
	}
}
