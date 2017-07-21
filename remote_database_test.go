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

	db, err := NewRemoteDatabase("mydb", ahr)
	assert.Nil(t, err)
	assert.Equal(t, "mydb", db.Name())
}

func TestFailsWhenLoadingNonExistingDatabase(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/mydb",
		httpmock.NewStringResponder(404, ""))

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	db, err := NewRemoteDatabase("mydb", ahr)
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

	db, err := NewRemoteDatabase("mydb", ahr)
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

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/mydb/_design/mydesigndoc",
		httpmock.NewStringResponder(200, `{
			"_id":"_design/mydesigndoc",
			"_rev":"4-aa671281619b1e9ae65e33995098e1d8",
			"language":"javascript",
			"views":{
				"by_user_id_and_created_at":{
					"map":"function(doc) {\n if ((doc['type'] == 'Authentication') && (doc['user_id'] != null) && (doc['created_at'] != null)) {\n emit([doc['user_id'], doc['created_at']], 1);\n }\n }",
					"reduce":"_sum"
				},
				"by_user_id_and_provider_and_uid_and_expires_at":{
					"map":"function(doc) {\n if ((doc['type'] == 'Authentication') && (doc['user_id'] != null) && (doc['provider'] != null) && (doc['uid'] != null) && (doc['expires_at'] != null)) {\n emit([doc['user_id'], doc['provider'], doc['uid'], doc['expires_at']], 1);\n }\n }",
					"reduce":"_sum"
				},
				"all":{"map":"function(doc) {\n if (doc['type'] == 'Authentication') {\n emit(doc._id, null);\n }\n }\n"}
			}
		}`))

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/mydb/_design/myseconddesigndoc",
		httpmock.NewStringResponder(200, `{
			"_id":"_design/myseconddesigndoc",
			"_rev":"2-aa671281619b1e9ae65e3afa334349855678",
			"language":"javascript",
			"views":{
				"by_user_id_and_provider_and_expires_at":{
					"map":"function(doc) {\n if ((doc['type'] == 'Authentication') && (doc['user_id'] != null) && (doc['provider'] != null) && (doc['expires_at'] != null)) {\n emit([doc['user_id'], doc['provider'], doc['expires_at']], 1);\n }\n }",
					"reduce":"_sum"
				},
				"all":{"map":"function(doc) {\n if (doc['type'] == 'Authentication') {\n emit(doc._id, null);\n }\n }\n"}
			}
		}`))

	docs, err := db.AllDesignDocs()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(docs))
	assert.Equal(t, "mydesigndoc", docs[0].Id())
	assert.Equal(t, "myseconddesigndoc", docs[1].Id())
}
