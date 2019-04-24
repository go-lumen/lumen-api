package server

import (
	"github.com/globalsign/mgo"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

// SetupMongoDatabase establishes the connexion with the mongo database
func (a *API) SetupMongoDatabase() (*mgo.Session, error) {
	session, err := mgo.Dial(a.Config.GetString("mongo_db_host"))
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	a.MongoDatabase = session.DB(a.Config.GetString("mongo_db_name"))

	return session, nil
}

// SetupPostgreDatabase establishes the connexion with the postgre database
func (a *API) SetupPostgreDatabase() (*pg.DB, error) {
	pgOptions := &pg.Options{
		//Addr:     a.Config.GetString("postgre_db_addr"),
		Database: a.Config.GetString("postgre_db_dbname"),
		User:     a.Config.GetString("postgre_db_user"),
	}
	db := pg.Connect(pgOptions)
	a.PostgreDatabase = db

	return db, nil
}
