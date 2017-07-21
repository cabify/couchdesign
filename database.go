package couchdesign

type Database interface {
	AllDesignDocs() ([]DesignDoc, error)
	Name() string
}
