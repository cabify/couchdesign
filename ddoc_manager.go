package couchdesign

import (
	"fmt"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/cabify/couchdesign/file_utils"
)

type DDocVersion map[Location]DesignDocsIndex
type DatabaseSummary map[string]DDocVersion

type Location uint8

const (
	LOCAL Location = iota
	REMOTE
)

type DDocManager struct {
	server   *Server
	localDir string
}

func NewDDocManager(ahr *AuthHttpRequester, localDir string) (*DDocManager, error) {
	var s *Server
	s, err := NewServer(ahr)
	if err != nil {
		return nil, err
	}

	if file_utils.MissingDir(localDir) {
		log.WithError(err).WithField("localDir", localDir).Error("Couldn't locate directory!")
		return nil, fmt.Errorf("Could not open directory %s", localDir)
	}

	return &DDocManager{
		server:   s,
		localDir: localDir,
	}, nil
}

func (d *DDocManager) Status() (map[string]DatabaseSummary, error) {
	log.WithField("server", d.server).Info("Loading remote databases...")
	dbList, err := d.server.AllDbs()
	if err != nil {
		log.WithError(err).Error("Could not get all databases!")
		return nil, err
	}

	dbs := make(map[string]DatabaseSummary, len(dbList))

	for _, db := range dbList {
		dbs[db.Name()] = make(map[string]DDocVersion, 0)
		if err := fetchAndIndexDesignDocs(db, dbs[db.Name()], REMOTE); err != nil {
			return nil, err
		}
	}

	log.WithField("directory", d.localDir).Info("Loading local databases...")
	localDbsList, err := ioutil.ReadDir(d.localDir)
	if err != nil {
		log.WithError(err).WithField("directory", d.localDir).Error("Could not read local directory!")
		return nil, err
	}

	for _, localDbName := range localDbsList {
		db, err := NewLocalDatabase(fmt.Sprintf("%s/%s", d.localDir, localDbName.Name()))
		if err != nil {
			log.WithError(err).WithField("directory", localDbName.Name()).Error("Could not create local database!")
			return nil, err
		}

		if err := fetchAndIndexDesignDocs(db, dbs[db.Name()], LOCAL); err != nil {
			return nil, err
		}
	}
	return dbs, nil
}

func fetchAndIndexDesignDocs(db Database, index map[string]DDocVersion, location Location) error {
	ddList, err := db.AllDesignDocs()
	if err != nil {
		log.WithError(err).WithField("database", db.Name()).Error("Could not get design docs!")
		return err
	}

	for _, dd := range ddList {
		if _, ok := index[dd.Id()]; !ok {
			index[dd.Id()] = make(map[Location]DesignDocsIndex)
		}
		index[dd.Id()][location] = dd
	}

	return nil
}
