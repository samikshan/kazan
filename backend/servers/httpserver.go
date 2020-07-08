package servers

import (
	"github.com/labstack/echo"
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

	sv := &HttpServer{
		E: echo.New(),
		H: handlers.New(userRepo),
	}

	sv.E.Logger.SetLevel(glog.DEBUG)

	sv.setupRoutes()

	return sv
}

func (sv *HttpServer) setupRoutes() {
	sv.E.POST("/signup", sv.H.Signup)
	sv.E.POST("/login", sv.H.Login)
	sv.E.POST("/token", sv.H.Token)
}
