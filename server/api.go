package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// Environment define on which env. we are currently running the server
type Environment = string

// Define environments
const (
	EnvLocal Environment = "local"
	EnvProd  Environment = "prod"
)

// Define configuration keys
const (
	ConfigDbType      = "db_type"
	ConfigHostAddress = "host_address"
)

// API structure that holds various necessary services
type API struct {
	Router          *gin.Engine
	Config          *viper.Viper
	MongoDatabase   *mongo.Database
	PostgreDatabase *gorm.DB
	MySQLDatabase   *sql.DB
	EmailSender     services.EmailSender
	TextSender      services.TextSender
}

// NewAPI allows to instantiate an API structure
func NewAPI() *API {
	return &API{Config: viper.New()}
}

// GetEnv returns the current environment defined in SAM_ENV
func (a *API) GetEnv() Environment {
	switch os.Getenv("SAM_ENV") {
	case EnvLocal:
		return EnvLocal
	case EnvProd:
		return EnvProd
	default:
		return EnvLocal
	}
}

// Setup bootstraps the sever
func (a *API) Setup() error {
	// Configuration setup
	err := a.SetupViper()
	utils.CheckErr(err)
	if err != nil {
		logrus.WithError(err).Errorln("cannot load configuration")
		return err
	}

	var router *gin.Engine

	if a.GetEnv() == EnvProd {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	} else {
		router = gin.Default()
	}
	a.Router = router

	router.Use(gin.Recovery())

	// Email sender setup
	a.EmailSender = services.NewEmailSender(a.Config)
	a.TextSender = services.NewTextSender(a.Config)

	// Database setup
	dbType := a.Config.GetString(ConfigDbType)
	logrus.Infoln("dbType:", dbType)
	switch dbType {
	case "mongo":
		_, err := a.SetupMongoDatabase()
		if err == nil {
			utils.Log(nil, "info", "SetupMongoDatabase OK")
		} else {
			utils.Log(nil, "err", "SetupMongoDatabase KO:", err)
		}
		utils.CheckErr(err)
		//defer session.Close()

		err = a.SetupMongoIndexes()
		if err == nil {
			utils.Log(nil, "info", "SetupMongoIndexes OK")
		} else {
			utils.Log(nil, "err", "SetupMongoIndexes KO:", err)
		}
		utils.CheckErr(err)

		// Seeds setup
		err = a.SetupMongoSeeds()
		utils.CheckErr(err)
		if err == nil {
			utils.Log(nil, "info", "SetupMongoSeeds OK")
		} else {
			utils.Log(nil, "err", "SetupMongoSeeds KO:", err)
		}
		utils.CheckErr(err)

	case "postgresql":
		panic("not supported")
	}

	// Router setup
	a.SetupRouter()
	logrus.Infoln("router started successfully")

	return err
}

// ScheduleCronTasks starts scheduled tasks
func (a *API) ScheduleCronTasks() error {
	c := cron.New()
	for _, task := range tasks {
		if _, err := c.AddFunc(task.CronSpec, func() {
			// cron task are ran as god
			ctx := store.NewGodContext(mongodb.New(new(gin.Context), a.MongoDatabase, a.Config.GetString("mongo_db_name")))
			task.Handler(ctx)
		}); err != nil {
			return errors.Wrap(err, "cannot start task: "+task.Name)
		}
		logrus.WithField("name", task.Name).Info("Started task successfully")
	}
	c.Start()
	return nil
}

// Run allows to run the server
func (a *API) Run() error {
	if err := a.ScheduleCronTasks(); err != nil {
		return errors.Wrap(err, "Cannot run server")
	}
	return a.Router.Run(a.Config.GetString(ConfigHostAddress))
}
