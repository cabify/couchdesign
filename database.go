package couchdesign

import (
	"fmt"

	"github.com/apex/log"
)

type Database struct {
	name string
}

func NewDatabase(name string, ahr *AuthHttpRequester) (*Database, error) {
	if err := ahr.Head(name); err != nil {
		return nil, err
	}

	return &Database{
		name: name,
	}, nil
}

func (d *Database) AllDesignDocs(ahr *AuthHttpRequester) ([]*DesignDoc, error) {
	log.WithField("database", d.name).Info("Fetching Design docs list...")
	var data struct {
		Rows []*DesignDoc `json:"rows"`
	}
	if err := ahr.Get(fmt.Sprintf("%s/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"", d.name), &data); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

func (d *Database) Name() string {
	return d.name
}
