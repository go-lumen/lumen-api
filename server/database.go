package server

import (
	"database/sql"
	"github.com/globalsign/mgo"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/go-pg/pg"
	_ "github.com/go-sql-driver/mysql" // For MySQL
	//_ "github.com/lib/pq"              // For PostgreSQL
	"github.com/sirupsen/logrus"
	//https://github.com/jinzhu/gorm
	//https://github.com/go-xorm/xorm
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	sql, _ := q.FormattedQuery()
	logrus.Debugf("SQL query: %s", sql)
}

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
func (a *API) SetupPostgreDatabase() (*pg.DB, error) {
	pgOptions := &pg.Options{
		Addr:     a.Config.GetString("postgres_db_addr"),
		Database: a.Config.GetString("postgres_db_name"),
		User:     a.Config.GetString("postgres_db_user"),
	}
	db := pg.Connect(pgOptions)

	db.AddQueryHook(dbLogger{})
	/*db, err := sql.Open("postgres", "user="+a.Config.GetString("postgres_db_user")+
		" dbName="+a.Config.GetString("postgres_db_dbname")+" sslmode=verify-full")
	utils.CheckErr(err)*/

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
