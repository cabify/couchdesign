package couchdesign

import "github.com/apex/log"

type Server struct {
	Version string `json:"version"`
	Uuid    string `json:"uuid"`
}

func NewServer(ahr *AuthHttpRequester) (*Server, error) {
	log.WithField("server", ahr.BaseUrl()).Info("Connecting to server...")
	var s Server
	if err := ahr.Get("", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Server) AllDbs(ahr *AuthHttpRequester) ([]*Database, error) {
	var dbNames []string
	log.WithField("server", ahr.BaseUrl()).Info("Listing databases at server...")
	if err := ahr.Get("_all_dbs", &dbNames); err != nil {
		return nil, err
	}

	dbs := make([]*Database, len(dbNames))
	for i, dbName := range dbNames {
		db, err := NewDatabase(dbName)
		if err != nil {
			return nil, err
		}
		dbs[i] = db
	}
	return dbs, nil
}
