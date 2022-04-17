package main

import (
	"github.com/go-lumen/lumen-api/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	api := server.NewAPI()
	err := api.Setup()

	if err != nil {
		logrus.WithError(err).Errorln("cannot setup api")
		os.Exit(1)
	}

	err = api.Run()
	if err != nil {
		logrus.WithError(err).Errorln("cannot run api")
		os.Exit(1)
	}
}
