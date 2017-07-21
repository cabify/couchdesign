package couchdesign

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistingLocalDatabaseIsLoaded(t *testing.T) {
	dir, err := ioutil.TempDir("", "test_db")
	assert.Nil(t, err)

	db, err := NewLocalDatabase(dir)
	assert.Nil(t, err)

	assert.True(t, strings.HasPrefix(db.Name(), "test_db"), "Database name is wrong!")

	os.Remove(dir)
}

func TestLoadDbFailsIfDirNotFound(t *testing.T) {
	db, err := NewLocalDatabase("/path/to/non/existing/dir")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestLocalDbLoadsDesignDocs(t *testing.T) {
	dir, err := ioutil.TempDir("", "test_db")
	assert.Nil(t, err)

	dir1, err := ioutil.TempDir(dir, "firstddoc")
	assert.Nil(t, err)
	writeYamlInto(dir1, "myfirstddoc", t)

	dir2, err := ioutil.TempDir(dir, "secondddoc")
	assert.Nil(t, err)
	writeYamlInto(dir2, "mysecondddoc", t)
	writeYamlInto(dir2, "mysecondddoc", t)

	_, err = ioutil.TempDir(dir, "emptyddoc")
	assert.Nil(t, err)

	db, err := NewLocalDatabase(dir)
	assert.Nil(t, err)

	ddocs, err := db.AllDesignDocs()
	assert.Nil(t, err)

	assert.Len(t, ddocs, 3)
	assert.True(t, strings.HasPrefix(ddocs[0].Id(), "myfirstddoc"), "Wrong ddoc id")
	assert.True(t, strings.HasPrefix(ddocs[1].Id(), "mysecondddoc"), "Wrong ddoc id")
	assert.True(t, strings.HasPrefix(ddocs[2].Id(), "mysecondddoc"), "Wrong ddoc id")
}

func writeYamlInto(dir, prefix string, t *testing.T) {
	file, err := ioutil.TempFile(dir, prefix)
	assert.Nil(t, err)

	data := "language: javascript"

	err = ioutil.WriteFile(file.Name(), []byte(data), os.ModeTemporary)
	assert.Nil(t, err)
}
