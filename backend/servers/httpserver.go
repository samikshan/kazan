package servers

import (
	"github.com/labstack/echo"
)

type HttpServer struct {
	E *echo.Echo
}

func NewHTTPSv() *HttpServer {
	return &HttpServer{
		E: echo.New(),
	}
}
