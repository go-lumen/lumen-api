package server

import "github.com/globalsign/mgo"

// SetupDatabase establishes the connexion with the database
func (a *API) SetupDatabase() (*mgo.Session, error) {
	session, err := mgo.Dial(a.Config.GetString("db_host"))
	if err != nil {
		return nil, err
	}

	a.MongoDatabase = session.DB(a.Config.GetString("db_name"))

	return session, nil
}
