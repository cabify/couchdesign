package couchdesign

import "fmt"

type Database struct {
	name string
}

func NewDatabase(name string) (*Database, error) {
	return &Database{
		name: name,
	}, nil
}

func (d *Database) AllDesignDocs(ahr *AuthHttpRequester) ([]*DesignDoc, error) {
	var data struct {
		results []*DesignDoc `json:"results"`
	}
	if err := ahr.Get(fmt.Sprintf("%s/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"&include_docs=true", d.name), &data); err != nil {
		return nil, err
	}
	return data.results, nil
}

func (d *Database) Name() string {
	return d.name
}
