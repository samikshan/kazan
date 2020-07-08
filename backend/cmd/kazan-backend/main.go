package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/samikshan/kazan/backend/servers"
)

func main() {
	sv := servers.NewHTTPSv()
	log.Fatal(sv.E.Start(":1323"))
}
