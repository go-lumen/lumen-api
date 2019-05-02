package server

import (
	"database/sql"
	"github.com/globalsign/mgo"
	"github.com/go-lumen/lumen-api/utils"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
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
func (a *API) SetupPostgreDatabase() (*sql.DB, error) {
	/*pgOptions := &pg.Options{
		//Addr:     a.Config.GetString("postgre_db_addr"),
		Database: a.Config.GetString("postgre_db_dbname"),
		User:     a.Config.GetString("postgre_db_user"),
	}
	db := pg.Connect(pgOptions)*/

	db, err := sql.Open("postgres", "user="+a.Config.GetString("postgre_db_user")+
		" dbName="+a.Config.GetString("postgre_db_dbname")+" sslmode=verify-full")
	utils.CheckErr(err)

	a.PostgreDatabase = db

	return db, nil
}

func (a *API) SetupMySQLDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", a.Config.GetString("mysql_db_user") + ":" + a.Config.GetString("mysql_db_password") +
		"@/" + a.Config.GetString("mysql_db_dbname"))
	utils.CheckErr(err)

	a.MySQLDatabase = db

	return db, nil
}