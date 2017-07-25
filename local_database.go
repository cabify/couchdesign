package couchdesign

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/apex/log"
	"github.com/cabify/couchdesign/file_utils"
)

type LocalDatabase struct {
	name, localDir string
}

func NewLocalDatabase(localDir string) (*LocalDatabase, error) {
	if file_utils.MissingDir(localDir) {
		log.WithField("localDir", localDir).Error("Couldn't locate directory!")
		return nil, fmt.Errorf("Could not open directory %s", localDir)
	}

	parts := strings.Split(localDir, "/")

	return &LocalDatabase{
		name:     parts[len(parts)-1],
		localDir: localDir,
	}, nil
}

func (d *LocalDatabase) AllDesignDocs() (map[string]DesignDocsIndex, error) {
	log.WithField("database", d.name).Info("Reading local design docs list...")
	ddList, err := ioutil.ReadDir(d.localDir)
	if err != nil {
		log.WithError(err).WithField("database", d.name).Error("Could not read directory!")
		return nil, err
	}

	dds := make(map[string]DesignDocsIndex, len(ddList))
	for _, ddName := range ddList {
		idx, err := NewLocalDesignDocsIndex(d.localDir, ddName.Name())
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"database": d.name, "design_doc": ddName.Name()}).
				Error("Could not create the local design docs index!")
			return nil, err
		}
		dds[ddName.Name()] = idx
	}

	return dds, nil
}

func (d *LocalDatabase) Name() string {
	return d.name
}
