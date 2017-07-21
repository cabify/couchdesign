package couchdesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestDesignDocIsLoaded(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	dd, err := NewRemoteDesignDoc("mydb", "mydesigndoc", ahr)
	assert.Nil(t, err)
	assert.Equal(t, "mydesigndoc", dd.Id())
	assert.Equal(t, "javascript", dd.Contents.Language)
	assert.Len(t, dd.Contents.Views, 3)
	assert.Contains(t, dd.Contents.Views, "by_user_id_and_created_at")
	assert.Contains(t, dd.Contents.Views, "by_user_id_and_provider_and_uid_and_expires_at")
	assert.Contains(t, dd.Contents.Views, "all")
}

func TestNonExistingDesignDocReturnsError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	dd, err := NewRemoteDesignDoc("mydb", "mydesigndoc", ahr)
	assert.Nil(t, dd)
	assert.Error(t, err)
}
