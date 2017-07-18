package couchdesign

import (
	"testing"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestExistingDatabaseIsLoaded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/mydb",
		httpmock.NewStringResponder(200, ""))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	db, err := NewDatabase("mydb", ahr)
	assert.Nil(t, err)
	assert.Equal(t, "mydb", db.Name())
}

func TestNonExistingDatabaseReturnsError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/mydb",
		httpmock.NewStringResponder(404, ""))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	db, err := NewDatabase("mydb", ahr)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDbLoadsAllDesignDocs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/mydb",
		httpmock.NewStringResponder(200, ""))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	db, err := NewDatabase("mydb", ahr)
	assert.Nil(t, err)

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/mydb/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"",
		httpmock.NewStringResponder(200, `{
			"total_rows":13,
			"offset":7,
			"rows":[
				{ "id":"_design/mydesigndoc", "key":"_design/mydesigndoc", "value":{ "rev":"2-509199e81180941112092f3fe3139643" } },
				{ "id":"_design/myseconddesigndoc", "key":"_design/myseconddesigndoc", "value":{ "rev":"32-556799e8167uj41112092f3fe3139098" } }
				]
			}`))

	docs, err := db.AllDesignDocs(ahr)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(docs))
	assert.Equal(t, "mydesigndoc", docs[0].Id())
	assert.Equal(t, "myseconddesigndoc", docs[1].Id())
}
