package couchdesign

import (
	"testing"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestNewServerConnectsToServer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/",
		httpmock.NewStringResponder(200, `{
			"couchdb": "Welcome",
			"uuid": "8cd8b27a02b73bbc3cecfaafe81df742",
			"version": "1.6.1",
			"vendor": {
				"version": "1.6.1",
				"name": "The Apache Software Foundation"
			}
		}`))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	server, err := NewServer(ahr)

	assert.Nil(t, err)
	assert.Equal(t, "1.6.1", server.Version)
	assert.Equal(t, "8cd8b27a02b73bbc3cecfaafe81df742", server.Uuid)
}

func TestNewServerFailsIfCannotConnect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	server, err := NewServer(ahr)
	assert.Nil(t, server)
	assert.Error(t, err, "An error should have been returned as the server cannot be contacted")
}

func TestListAllDbs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/",
		httpmock.NewStringResponder(200, `{
			"couchdb": "Welcome",
			"uuid": "8cd8b27a02b73bbc3cecfaafe81df742",
			"version": "1.6.1",
			"vendor": {
				"version": "1.6.1",
				"name": "The Apache Software Foundation"
			}
		}`))

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/_all_dbs",
		httpmock.NewStringResponder(200, `["db1", "db2"]`))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	server, err := NewServer(ahr)
	assert.Nil(t, err)

	dbs, err := server.AllDbs(ahr)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(dbs))
	assert.Equal(t, "db1", dbs[0].Name())
	assert.Equal(t, "db2", dbs[1].Name())
}
