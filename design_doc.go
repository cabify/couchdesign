package couchdesign

import "strings"

type DesignDoc struct {
	FullId string `json:"id"`
}

func (d *DesignDoc) Id() string {
	return strings.SplitN(d.FullId, "/", 2)[1]
}
