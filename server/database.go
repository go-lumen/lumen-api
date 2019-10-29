package server

import (
	"database/sql"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

type dbLogger struct{}

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

// SetupPostgreDatabase establishes the connexion with the PostgreSQL database
func (a *API) SetupPostgreDatabase() (*gorm.DB, error) {
	/*connStr := string("host=" + a.Config.GetString("postgres_db_addr") + " port=" + a.Config.GetString("postgres_db_port") +
		" user=" + a.Config.GetString("postgres_db_user") + " dbname=" + a.Config.GetString("postgres_db_name") +
		" password=" + a.Config.GetString("postgres_db_password"))*/
	db, err := gorm.Open("postgres", "host=localhost port=5432 sslmode=disable user=adrien password=litfsoh:PQ dbname=lumen")
	fmt.Println("error:", err)
	//utils.CheckErr(err)
	defer db.Close()

	a.PostgreDatabase = db

	return db, nil
}

// SetupMySQLDatabase establishes the connexion with the MySQL database
func (a *API) SetupMySQLDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", a.Config.GetString("mysql_db_user")+":"+a.Config.GetString("mysql_db_password")+
		"@/"+a.Config.GetString("mysql_db_name"))
	utils.CheckErr(err)

	a.MySQLDatabase = db

	return db, nil
}
