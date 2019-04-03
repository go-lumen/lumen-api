package server

import "github.com/globalsign/mgo"

// SetupMongoDatabase establishes the connexion with the database
func (a *API) SetupMongoDatabase() (*mgo.Session, error) {
	session, err := mgo.Dial(a.Config.GetString("db_host"))
	if err != nil {
		return nil, err
	}

	a.MongoDatabase = session.DB(a.Config.GetString("db_name"))

	return session, nil
}
