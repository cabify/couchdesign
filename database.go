package couchdesign

type Database interface {
	AllDesignDocs() (map[string]DesignDocsIndex, error)
	Name() string
}
