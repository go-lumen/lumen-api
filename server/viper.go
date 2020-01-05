package server

import (
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
)

// SetupViper according to env var that selects related conf file
func (a *API) SetupViper() error {
	filename := ".env"
	logrus.Infoln("Using env:" + os.Getenv("LUMEN_ENV"))
	switch os.Getenv("LUMEN_ENV") {
	case "testing":
		filename = "../.env.testing"
	case "prod":
		filename = ".env.prod"
	}

	err := godotenv.Overload(filename)
	utils.CheckErr(err)

	a.Config.SetEnvPrefix("lumen")
	a.Config.AutomaticEnv()

	return nil
}
