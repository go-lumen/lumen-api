package services

import (
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// KPITask processes kpis
func KPITask(ctx *store.Context) {
	now := time.Now()
	fromT := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)
	toT := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.UTC)

	logrus.Info("Daily kpis at : ", now.Format("2006.01.02 15:04:05"), "from:", fromT.Format("2006.01.02 15:04:05"), "to:", toT.Format("2006.01.02 15:04:05"), "(from:", fromT.Unix(), "to:", toT.Unix(), ")")

	dbUsers, _ := models.GetUsers(ctx, bson.M{})

	for _, user := range dbUsers {
		userCreationTS := utils.MongoIDToTimestamp(user.ID.Hex()).Unix()
		if (userCreationTS > fromT.Unix()) && (userCreationTS < toT.Unix()) {
			utils.Log(ctx.C, "info", "New user registered:", user.Email)
		}
	}
}
