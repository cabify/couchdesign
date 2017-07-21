package couchdesign

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalDesignDocumentsAreLoaded(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test_dd")
	assert.Nil(t, err)

	data := `
language: javascript
views:
  by_name:
    map: |
      function(doc) {
        emit(doc.name, doc)
      }
    reduce: _sum
  by_id_card:
    map: |
      function(doc) {
        emit(doc.card, 1)
      }
`
	err = ioutil.WriteFile(file.Name(), []byte(data), os.ModeTemporary)
	assert.Nil(t, err)

	dd, err := NewLocalDesignDoc(file.Name())
	assert.Nil(t, err)

	assert.True(t, strings.HasPrefix(dd.Id(), "test_dd"), "Wrong DD Id!")
	assert.Equal(t, "javascript", dd.Contents.Language)
	assert.Len(t, dd.Contents.Views, 2)
	assert.Contains(t, dd.Contents.Views, "by_name")
	assert.Contains(t, dd.Contents.Views, "by_id_card")

	os.Remove(file.Name())
}

func TestLoadFailsIfFileIsNotFound(t *testing.T) {
	dd, err := NewLocalDesignDoc("/not/found/file")
	assert.Error(t, err)
	assert.Nil(t, dd)
}

func TestLoadFailsIfMalformedYAML(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "test_dd")
	assert.Nil(t, err)

	data := "invalid yaml"

	err = ioutil.WriteFile(file.Name(), []byte(data), os.ModeTemporary)
	assert.Nil(t, err)

	dd, err := NewLocalDesignDoc(file.Name())
	assert.Error(t, err)
	assert.Nil(t, dd)

	os.Remove(file.Name())
}
