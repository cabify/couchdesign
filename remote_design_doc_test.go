package couchdesign

import (
	"testing"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cabify/couchdesign/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestDesignDocIsLoaded(t *testing.T) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	ddName := test_utils.RandDesignDocName()

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://127.0.0.1:5984/couchdesigntest/_design/%s", ddName),
		strings.NewReader(test_utils.Example1Json))
	assert.Nil(t, err)

	req.SetBasicAuth("admin", "pass")

	res, err := httpClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, res.StatusCode)

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	dd, err := NewRemoteDesignDoc("couchdesigntest", ddName, ahr)
	assert.Nil(t, err)
	assert.Equal(t, ddName, dd.Id())
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

func TestNonExistingDesignDocReturnsError(t *testing.T) {
	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	dd, err := NewRemoteDesignDoc("couchdesigntest", "nonexistentdesigndoc", ahr)
	assert.Nil(t, dd)
	assert.Error(t, err)
}

func TestRemoteDesignDocMD5IsGeneratedProperly(t *testing.T) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("PUT", "http://127.0.0.1:5984/couchdesigntest/_design/test_dd",
		strings.NewReader(test_utils.Example1Json))
	assert.Nil(t, err)

	req.SetBasicAuth("admin", "pass")

	res, err := httpClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 201, res.StatusCode)

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	dd, err := NewRemoteDesignDoc("couchdesigntest", "test_dd", ahr)
	assert.Nil(t, err)
	assert.Equal(t, "test_dd", dd.Id())
	assert.Equal(t, "4b30ab1cb9d313df921c1073a3b15f70", dd.Md5())

	req, err = http.NewRequest("DELETE", fmt.Sprintf("http://127.0.0.1:5984/couchdesigntest/_design/test_dd?rev=%s", dd.Contents.Rev), nil)
	assert.Nil(t, err)
	req.SetBasicAuth("admin", "pass")

	res, err = httpClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
