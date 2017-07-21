package couchdesign

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

type DesignDocContents struct {
	FullId            string              `json:"_id"`
	Language          string              `json:"language"`
	Options           map[string]bool     `json:"options"`
	Filters           map[string]Function `json:"filters"`
	Lists             map[string]Function `json:"lists"`
	Shows             map[string]Function `json:"shows"`
	Updates           map[string]Function `json:"updates"`
	ValidateDocUpdate Function            `json:"validate_doc_update"`
	Views             map[string]*View    `json:"views"`
} // TODO add rewrites key

type DesignDoc interface {
	Id() string
	Md5() string
}

func DesignDocMd5(d DesignDocContents) string {
	h := md5.New()

	io.WriteString(h, d.FullId)

	io.WriteString(h, "Language")
	io.WriteString(h, d.Language)

	io.WriteString(h, "Options")
	for key, opt := range d.Options {
		io.WriteString(h, key)
		io.WriteString(h, strconv.FormatBool(opt))
	}

	io.WriteString(h, "Filters")
	for key, filter := range d.Filters {
		io.WriteString(h, key)
		io.WriteString(h, string(filter))
	}

	io.WriteString(h, "Lists")
	for key, list := range d.Lists {
		io.WriteString(h, key)
		io.WriteString(h, string(list))
	}

	io.WriteString(h, "Shows")
	for key, show := range d.Shows {
		io.WriteString(h, key)
		io.WriteString(h, string(show))
	}

	io.WriteString(h, "Updates")
	for key, update := range d.Updates {
		io.WriteString(h, key)
		io.WriteString(h, string(update))
	}

	io.WriteString(h, "ValidateDocUpdate")
	io.WriteString(h, string(d.ValidateDocUpdate))

	io.WriteString(h, "Views")
	for key, view := range d.Views {
		io.WriteString(h, key)
		if view.Map != nil {
			io.WriteString(h, string(*view.Map))
		}
		if view.Reduce != nil {
			io.WriteString(h, string(*view.Reduce))
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
