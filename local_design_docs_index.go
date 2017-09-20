package couchdesign

import (
	"fmt"
	"io/ioutil"

	"strings"

	"strconv"

	"github.com/apex/log"
)

type LocalDesignDocsIndex struct {
	id       string
	versions map[int64]DesignDoc
}

func NewLocalDesignDocsIndex(basedir, ddName string) (DesignDocsIndex, error) {
	ddVersions, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", basedir, ddName))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"design_doc": ddName}).
			Error("Could not read dd's dir")
		return nil, err
	}

	ldi := &LocalDesignDocsIndex{
		id:       ddName,
		versions: make(map[int64]DesignDoc, len(ddVersions)),
	}

	for _, ddVersionName := range ddVersions {
		dd, err := NewLocalDesignDoc(fmt.Sprintf("%s/%s", ddName, ddVersionName.Name()))
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"design_doc": ddName, "version": ddVersionName}).
				Error("Could not create dd version!")
			return nil, err
		}
		version, err := versionFromDDName(ddVersionName.Name())
		if err != nil {
			log.WithError(err).
				WithFields(log.Fields{"design_doc": ddName, "version_name": ddVersionName}).
				Error("Could not parse version int!")
			return nil, err
		}
		ldi.versions[version] = dd
	}
	return ldi, nil
}

func versionFromDDName(name string) (int64, error) {
	parts := strings.Split(name, "_")
	return strconv.ParseInt(parts[len(parts)-1], 10, 32)
}

func (l *LocalDesignDocsIndex) Id() string {
	return l.id
}

func (l *LocalDesignDocsIndex) Versions() map[int64]DesignDoc {
	return l.versions
}
