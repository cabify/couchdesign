package couchdesign

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestDDocManagerIsLoaded(t *testing.T) {
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

	mgr, err := NewDDocManager(ahr, os.TempDir())
	assert.Nil(t, err)
	assert.NotNil(t, mgr)
}

func TestGetStatusBuildsProperIndex(t *testing.T) {
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

	dir, err := ioutil.TempDir("", "testdata")
	assert.Nil(t, err)
	defer os.Remove(dir)

	ahr, err := NewAuthHttpRequester("admin", "pass", "127.0.0.1")
	assert.Nil(t, err)

	mgr, err := NewDDocManager(ahr, dir)
	assert.Nil(t, err)

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/_all_dbs",
		httpmock.NewStringResponder(200, `["db1", "db2"]`))

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/db1",
		httpmock.NewStringResponder(200, ""))

	db1Dir := fmt.Sprintf("%s/%s", dir, "db1")
	err = os.Mkdir(db1Dir, os.ModePerm)
	assert.Nil(t, err)

	httpmock.RegisterResponder("HEAD", "http://127.0.0.1:5984/db2",
		httpmock.NewStringResponder(200, ""))

	db2Dir := fmt.Sprintf("%s/%s", dir, "db2")
	err = os.Mkdir(db2Dir, os.ModePerm)
	assert.Nil(t, err)

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/db1/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"",
		httpmock.NewStringResponder(200, `{
			"total_rows":13,
			"offset":7,
			"rows":[
				{ "id":"_design/mydesigndoc", "key":"_design/mydesigndoc", "value":{ "rev":"2-509199e81180941112092f3fe3139643" } },
				{ "id":"_design/myseconddesigndoc", "key":"_design/myseconddesigndoc", "value":{ "rev":"32-556799e8167uj41112092f3fe3139098" } }
				]
			}`))

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/db1/_design/mydesigndoc",
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

	myddDir := fmt.Sprintf("%s/%s", db1Dir, "mydesigndoc")
	err = os.Mkdir(myddDir, os.ModePerm)
	assert.Nil(t, err)

	data := `
language: javascript,
views:
  by_user_id_and_created_at:
    map: |
      function(doc) {
        if ((doc['type'] == 'Authentication') && (doc['user_id'] != null) && (doc['created_at'] != null)) {
          emit([doc['user_id'], doc['created_at']], 1);
        }
      }
    reduce: _sum
  by_user_id_and_provider_and_expires_at:
    map: |
      function(doc) {
        if ((doc['type'] == 'Authentication') && (doc['user_id'] != null) && (doc['provider'] != null) && (doc['uid'] != null) && (doc['expires_at'] != null)) {
          emit([doc['user_id'], doc['provider'], doc['uid'], doc['expires_at']], 1);
        }
      }
    reduce: _sum
  all:
    map: |
      function(doc) {
        if (doc['type'] == 'Authentication') {
          emit(doc._id, null);
        }
      }
`

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", myddDir, "mydesigndoc_1"), []byte(data), os.ModePerm)
	assert.Nil(t, err)

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/db1/_design/myseconddesigndoc",
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

	mysddDir := fmt.Sprintf("%s/%s", dir, "myseconddesigndoc")
	err = os.Mkdir(mysddDir, os.ModePerm)
	assert.Nil(t, err)

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/db2/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"",
		httpmock.NewStringResponder(200, `{
			"total_rows":13,
			"offset":7,
			"rows":[
				{ "id":"_design/greatdesigndoc", "key":"_design/greatdesigndoc", "value":{ "rev":"2-509199e81180941112092f3fe3139643" } }
				]
			}`))

	httpmock.RegisterResponder("GET", "http://127.0.0.1:5984/db2/_design/greatdesigndoc",
		httpmock.NewStringResponder(200, `{
			"_id":"_design/greatdesigndoc",
			"_rev":"2-aa671281619b1e9ae65e3afa334349855678",
			"language":"javascript",
			"views":{
				"by_name":{
					"map":"function(doc) {\n emit([doc['user_id'], doc['provider'], doc['expires_at']], 1);\n }",
					"reduce":"_sum"
				}
			}
		}`))

	index, err := mgr.Status()
	assert.Nil(t, err)
	assert.Len(t, index, 2)

	assert.Len(t, index["db1"], 2)
	assert.NotNil(t, index["db1"]["mydesigndoc"][REMOTE])
	assert.NotNil(t, index["db1"]["mydesigndoc"][LOCAL])

	assert.NotNil(t, index["db1"]["myseconddesigndoc"][REMOTE])
	assert.NotNil(t, index["db1"]["myseconddesigndoc"][LOCAL])

	assert.Len(t, index["db2"], 1)
	assert.NotNil(t, index["db2"]["greatdesigndoc"][REMOTE])
	assert.NotNil(t, index["db2"]["greatdesigndoc"][LOCAL])
}
