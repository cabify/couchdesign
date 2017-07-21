package couchdesign

import (
	"strings"

	"github.com/apex/log"
)

type Server struct {
	Version string `json:"version"`
	Uuid    string `json:"uuid"`
	ahr     *AuthHttpRequester
}

func NewServer(ahr *AuthHttpRequester) (*Server, error) {
	log.WithField("server", ahr.BaseUrl()).Info("Connecting to server...")
	s := &Server{
		ahr: ahr,
	}
	if err := ahr.Get("", s); err != nil {
		log.WithError(err).WithField("server", ahr.baseUrl).Error("Couldn't connect with server")
		return nil, err
	}
	return s, nil
}

func (s *Server) AllDbs() ([]Database, error) {
	var dbNames []string
	log.WithField("server", s.ahr.BaseUrl()).Info("Listing databases at server...")
	if err := s.ahr.Get("_all_dbs", &dbNames); err != nil {
		return nil, err
	}

	dbs := make([]Database, 0)
	for _, dbName := range dbNames {
		if strings.HasPrefix(dbName, "_") {
			continue
		}
		db, err := NewRemoteDatabase(dbName, s.ahr)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}
