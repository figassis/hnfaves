package api

import (
	"github.com/figassis/covidfree/pkg/utl/config"
	"github.com/figassis/covidfree/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) (err error) {

	e := server.New()

	e.GET("/:user", getFeed)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return
}
