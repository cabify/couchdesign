package couchdesign

import (
	"crypto/md5"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

type DesignDocContents struct {
	FullId            string              `json:"_id"`
	Rev               string              `json:"_rev"`
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

func DesignDocMd5(id string, d DesignDocContents) string {
	h := md5.New()

	io.WriteString(h, id)
	fmt.Println(d.FullId)

	io.WriteString(h, "Language")
	io.WriteString(h, d.Language)

	io.WriteString(h, "Options")
	keys := sortedKeys(reflect.ValueOf(d.Options))
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, strconv.FormatBool(d.Options[key]))
	}

	io.WriteString(h, "Filters")
	keys = sortedKeys(reflect.ValueOf(d.Filters))
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, string(d.Filters[key]))
	}

	io.WriteString(h, "Lists")
	keys = sortedKeys(reflect.ValueOf(d.Lists))
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, string(d.Lists[key]))
	}

	io.WriteString(h, "Shows")
	keys = sortedKeys(reflect.ValueOf(d.Shows))
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, string(d.Shows[key]))
	}

	io.WriteString(h, "Updates")
	keys = sortedKeys(reflect.ValueOf(d.Updates))
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, string(d.Updates[key]))
	}

	io.WriteString(h, "ValidateDocUpdate")
	io.WriteString(h, string(d.ValidateDocUpdate))

	io.WriteString(h, "Views")
	keys = sortedKeys(reflect.ValueOf(d.Views))
	for _, key := range keys {
		io.WriteString(h, key)
		view := d.Views[key]
		if view.Map != nil {
			io.WriteString(h, string(*view.Map))
		}
		if view.Reduce != nil {
			io.WriteString(h, string(*view.Reduce))
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func sortedKeys(data reflect.Value) []string {
	values := data.MapKeys()
	keys := make([]string, len(values))
	i := 0
	for _, k := range values {
		keys[i] = k.String()
		i++
	}

	sort.Strings(keys)
	return keys
}
