package couchdesign

//func TestNewLocalDesignDocsIndexLoadsVersionsProperly(t *testing.T) {
//	dir, err := ioutil.TempDir("", "testdata")
//	assert.Nil(t, err)
//	defer os.Remove(dir)
//
//	ddDir, err := test_utils.AddDirectoryToDir(dir, "mydesigndoc")
//	assert.Nil(t, err)
//
//	_, err = test_utils.WriteAnExample1DesignDocIntoFile(ddDir, "mydesigndoc_1")
//	assert.Nil(t, err)
//
//	_, err = test_utils.WriteAnExample1DesignDocIntoFile(ddDir, "mydesigndoc_2")
//	assert.Nil(t, err)
//
//	_, err = test_utils.WriteAnExample1DesignDocIntoFile(ddDir, "mydesigndoc_3")
//	assert.Nil(t, err)
//
//	index, err := NewLocalDesignDocsIndex(dir, "mydesigndoc")
//	assert.Nil(t, err)
//	assert.Equal(t, "mydesigndoc", index.Id())
//	assert.Len(t, index.Versions(), 3)
//}
