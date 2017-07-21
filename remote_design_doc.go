package couchdesign

import (
	"fmt"
	"strings"

	"github.com/apex/log"
)

type RemoteDesignDoc struct {
	Contents DesignDocContents
}

func NewRemoteDesignDoc(db, name string, ahr *AuthHttpRequester) (*RemoteDesignDoc, error) {
	log.WithFields(log.Fields{"database": db, "design_doc": name}).Info("Loading design doc...")
	dd := &RemoteDesignDoc{}
	if err := ahr.Get(fmt.Sprintf("%s/_design/%s", db, name), &dd.Contents); err != nil {
		return nil, err
	}
	return dd, nil
}

func (d *RemoteDesignDoc) Id() string {
	return designDocNameFromId(d.Contents.FullId)
}

func designDocNameFromId(id string) string {
	return strings.SplitN(id, "/", 2)[1]
}

func (d *RemoteDesignDoc) Md5() string {
	return DesignDocMd5(d.Contents)
}
