package test_utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var example1Yaml = `
language: javascript
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

var invalidYaml = `this, is, an invalid yaml`

func WriteAnExample1DesignDocIntoFile(t *testing.T, dir, file string) string {
	name := concat(dir, file)
	err := ioutil.WriteFile(name, []byte(example1Yaml), os.ModePerm)
	assert.Nil(t, err)
	return name
}

func WriteAnInvalidDesignDocIntoFile(t *testing.T, dir, file string) string {
	name := concat(dir, file)
	err := ioutil.WriteFile(name, []byte(invalidYaml), os.ModePerm)
	assert.Nil(t, err)
	return name
}

func AddDirectoryToDir(t *testing.T, parent, child string) string {
	name := concat(parent, child)
	err := os.Mkdir(name, os.ModePerm)
	assert.Nil(t, err)
	return name
}

func concat(str1, str2 string) string {
	return fmt.Sprintf("%s/%s", str1, str2)
}
