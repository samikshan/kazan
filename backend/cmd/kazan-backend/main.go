package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/samikshan/kazan/backend"
	"github.com/samikshan/kazan/backend/servers"
)

func main() {
	serverURL := fmt.Sprintf(":%d", backend.Cfg.HTTPPort)

	sv := servers.NewHTTPSv()
	log.Fatal(sv.E.Start(serverURL))
}
