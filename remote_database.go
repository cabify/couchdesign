package couchdesign

import (
	"github.com/apex/log"
)

type RemoteDatabase struct {
	name string
	ahr  *AuthHttpRequester
}

func NewRemoteDatabase(name string, ahr *AuthHttpRequester) (*RemoteDatabase, error) {
	db := &RemoteDatabase{
		name: name,
		ahr:  ahr,
	}

	if err := ahr.Head(name); err != nil {
		log.WithError(err).WithField("database", name).Error("Could not locate remote database!")
		return nil, err
	}

	return db, nil
}

func (d *RemoteDatabase) AllDesignDocs() (map[string]DesignDocsIndex, error) {
	log.WithField("database", d.name).Info("Fetching Design docs list...")
	//var data struct {
	//	Rows []struct {
	//		Id string `json:"id"`
	//	} `json:"rows"`
	//}
	//if err := d.ahr.Get(fmt.Sprintf("%s/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"", d.name), &data); err != nil {
	//	return nil, err
	//}
	//
	//for i, ddInfo := range data.Rows {
	//	dd, err := NewRemoteDesignDoc(d.name, designDocNameFromId(ddInfo.Id), d.ahr)
	//	if err != nil {
	//		return nil, err
	//	}
	//	dds[i] = dd
	//}

	return nil, nil
}

func (d *RemoteDatabase) Name() string {
	return d.name
}
