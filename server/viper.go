package server

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
)

// SetupViper according to env var that selects related conf file
func (a *API) SetupViper() error {
	filename := ".env"
	logrus.Infoln("Using env:" + os.Getenv("IOTHINGS_ENV"))
	switch os.Getenv("IOTHINGS_ENV") {
	case "testing":
		filename = "../.env.testing"
	case "prod":
		filename = ".env.prod"
	}

	godotenv.Overload(filename)

	a.Config.SetEnvPrefix("iothings")
	a.Config.AutomaticEnv()

	return nil
}
