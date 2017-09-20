package couchdesign

type DesignDocsIndex interface {
	Id() string
	Versions() map[int64]DesignDoc
}
