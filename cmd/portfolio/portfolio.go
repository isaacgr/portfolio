package main

import (
	"os"

	"github.com/isaacgr/portfolio/internal/pkg/logging"
	"github.com/isaacgr/portfolio/internal/web"
)

var log = logging.GetLogger("main", false)

func main() {
	c, err := web.NewConfiguration()
	if err != nil {
		log.Error("Unable to parse configuration. Exiting.")
		os.Exit(1)
	}
	server := web.NewWebServer(c)
	server.Start()
}
