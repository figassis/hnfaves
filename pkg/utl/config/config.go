package config

import (
	"fmt"
	"os"

	"github.com/figassis/hnfaves/pkg/utl/util"
	"github.com/figassis/hnfaves/pkg/utl/zaplog"
	"github.com/joho/godotenv"
)

var (
	environment = []string{"APP_PORT", "REQUESTS_PER_SECOND"}
)

const (
	version = "0.0.1"
)

func Version() string {
	return version
}

// Load returns Configuration struct
func Load() (cfg *Configuration, err error) {
	zaplog.Initialize(version)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if err = validateEnvironment(); err != nil {
		return
	}

	port := "80"
	if os.Getenv("APP_PORT") != "" {
		port = os.Getenv("APP_PORT")
	}

	cfg = &Configuration{
		Version: version,
		Server: &Server{
			Port:         port,
			Debug:        os.Getenv("ENVIRONMENT") != "production",
			ReadTimeout:  30, //seconds
			WriteTimeout: 30, //seconds
		},
	}

	//Initialize Cache
	if err = zaplog.ZLog(util.New()); err != nil {
		return
	}

	return
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Version string
	Server  *Server
}

// Server holds data necessery for server configuration
type Server struct {
	Port              string
	Debug             bool
	ReadTimeout       int
	WriteTimeout      int
	RequestsPerSecond int
}

func validateEnvironment() error {
	//If we have an .env fle, use that, otherwise, use existing environment
	if err := util.CheckPath(".env"); err == nil {
		if err = zaplog.ZLog(godotenv.Load()); err != nil {
			return zaplog.ZLog("Could not load env file")
		}
	}

	for _, env := range environment {
		if os.Getenv(env) == "" {
			return fmt.Errorf("Environment variable %s is required", env)
		}
	}
	return nil
}
