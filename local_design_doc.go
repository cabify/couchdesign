package couchdesign

import (
	"io/ioutil"
	"strings"

	"github.com/apex/log"

	yaml "gopkg.in/yaml.v2"
)

type LocalDesignDoc struct {
	Contents DesignDocContents
}

func NewLocalDesignDoc(fileName string) (*LocalDesignDoc, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithError(err).WithField("file", fileName).Error("Could not read file!")
		return nil, err
	}

	ldd := &LocalDesignDoc{}
	if err := yaml.Unmarshal(data, &ldd.Contents); err != nil {
		log.WithError(err).WithField("file", fileName).Error("Could not parse YAML file!")
		return nil, err
	}
	ldd.Contents.FullId = fileName

	return ldd, nil
}

func (d *LocalDesignDoc) Id() string {
	parts := strings.Split(d.Contents.FullId, "/")
	return parts[len(parts)-1]
}

func (d *LocalDesignDoc) Md5() string {
	return DesignDocMd5(d.Contents)
}
