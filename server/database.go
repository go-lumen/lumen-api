package server

import (
	"database/sql"
	"github.com/globalsign/mgo"
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

func (a *API) SetupPostgreDatabase() (*sql.DB, error) {
	connStr := "user=" + a.Config.GetString("postgre_db_user") + " dbname=" + a.Config.GetString("postgre_db_dbname") + " sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return db, nil
}
