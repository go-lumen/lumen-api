package server

import (
	"fmt"
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
	if err != nil {
		fmt.Println("godotenv error:", err)
	} else {
		logrus.Infoln("Godotenv OK")
	}

	a.Config.SetEnvPrefix("lumen")
	a.Config.AutomaticEnv()

	return nil
}
