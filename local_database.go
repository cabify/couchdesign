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

func (d *LocalDatabase) AllDesignDocs() ([]DesignDoc, error) {
	log.WithField("database", d.name).Info("Reading local design docs list...")
	ddList, err := ioutil.ReadDir(d.localDir)
	if err != nil {
		log.WithError(err).WithField("database", d.name).Error("Could not read directory!")
		return nil, err
	}

	dds := make([]DesignDoc, 0)
	for _, ddName := range ddList {
		ddBaseDir := fmt.Sprintf("%s/%s", d.localDir, ddName.Name())
		ddVersions, err := ioutil.ReadDir(ddBaseDir)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"database": d.name, "design_doc": ddName}).Error("Could not read dd's dir")
			return nil, err
		}
		for _, ddVersionName := range ddVersions {
			dd, err := NewLocalDesignDoc(fmt.Sprintf("%s/%s", ddBaseDir, ddVersionName.Name()))
			if err != nil {
				log.WithError(err).WithFields(log.Fields{"database": d.name, "design_doc": ddName, "version": ddVersionName}).Error("Could not create dd version!")
				return nil, err
			}
			dds = append(dds, dd)
		}
	}

	return dds, nil
}

func (d *LocalDatabase) Name() string {
	return d.name
}
