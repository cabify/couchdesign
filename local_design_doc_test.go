package couchdesign

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cabify/couchdesign/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestLocalDesignDocumentsAreLoaded(t *testing.T) {
	dir, err := ioutil.TempDir("", "testdata")
	assert.Nil(t, err)
	defer os.Remove(dir)

	file := test_utils.WriteAnExample1DesignDocIntoFile(t, dir, "test_dd.yaml")

	dd, err := NewLocalDesignDoc(file)
	assert.Nil(t, err)

	assert.Equal(t, "test_dd", dd.Id())
	assert.Equal(t, "javascript", dd.Contents.Language)
	assert.Len(t, dd.Contents.Views, 3)
	assert.Contains(t, dd.Contents.Views, "by_user_id_and_created_at")
	assert.NotNil(t, dd.Contents.Views["by_user_id_and_created_at"].Map)
	assert.NotNil(t, dd.Contents.Views["by_user_id_and_created_at"].Reduce)

	assert.Contains(t, dd.Contents.Views, "by_user_id_and_provider_and_expires_at")
	assert.NotNil(t, dd.Contents.Views["by_user_id_and_provider_and_expires_at"].Map)
	assert.NotNil(t, dd.Contents.Views["by_user_id_and_provider_and_expires_at"].Reduce)

	assert.Contains(t, dd.Contents.Views, "all")
	assert.NotNil(t, dd.Contents.Views["all"].Map)
	assert.Nil(t, dd.Contents.Views["all"].Reduce)
}

func TestLoadFailsIfFileIsNotFound(t *testing.T) {
	dd, err := NewLocalDesignDoc("/not/found/file")
	assert.Error(t, err)
	assert.Nil(t, dd)
}

func TestLoadFailsIfMalformedYAML(t *testing.T) {
	dir, err := ioutil.TempDir("", "testdata")
	assert.Nil(t, err)
	defer os.Remove(dir)

	file := test_utils.WriteAnInvalidDesignDocIntoFile(t, dir, "test_dd.yaml")

	dd, err := NewLocalDesignDoc(file)
	assert.Error(t, err)
	assert.Nil(t, dd)
}

func TestRequiresYamlExtensionFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "testdata")
	assert.Nil(t, err)
	defer os.Remove(dir)

	file := test_utils.WriteAnExample1DesignDocIntoFile(t, dir, "test_dd")

	dd, err := NewLocalDesignDoc(file)
	assert.Error(t, err)
	assert.Nil(t, dd)
}
