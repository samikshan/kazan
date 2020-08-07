package servers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	glog "github.com/labstack/gommon/log"
	"github.com/samikshan/kazan/backend/db"
	"github.com/samikshan/kazan/backend/handlers"
	"github.com/samikshan/kazan/backend/repositories"
)

type HttpServer struct {
	E *echo.Echo
	H *handlers.Handler
}

func NewHTTPSv() *HttpServer {
	dbConn := db.New()
	db.AutoMigrate(dbConn)

	userRepo := repositories.NewUserRepo(dbConn)
	authRepo := repositories.NewAuthRepo(dbConn)
	trackRepo := repositories.NewTrackRepo(dbConn)
	insRepo := repositories.NewInstrumentRepo(dbConn)

	sv := &HttpServer{
		E: echo.New(),
		H: handlers.New(userRepo, authRepo, trackRepo, insRepo),
	}

	sv.E.Use(middleware.Logger())
	sv.E.Use(middleware.Recover())
	sv.E.Use(middleware.CORS())
	sv.E.Logger.SetLevel(glog.DEBUG)

	sv.setupRoutes()

	return sv
}

func (sv *HttpServer) setupRoutes() {
	sv.E.POST("/authentication", sv.H.SetNewAuth)
	sv.E.GET("/authentication", sv.H.GetAuth)
	sv.E.GET("/user", sv.H.GetSender)
	sv.E.POST("/user", sv.H.CreateNewUser)
	sv.E.PUT("/user/:id", sv.H.UpdateUser)
	sv.E.GET("/feed", sv.H.GetUserFeed)
	sv.E.POST("/tracks", sv.H.NewTrack)
	sv.E.GET("/profiles/:id", sv.H.GetProfile)
}
