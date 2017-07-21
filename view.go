package couchdesign

type View struct {
	Map    *Function `json:"map"`
	Reduce *Function `json:"reduce"`
}
