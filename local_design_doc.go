package couchdesign

import (
	"io/ioutil"
	"strings"

	"github.com/apex/log"

	"github.com/juju/errors"
	yaml "gopkg.in/yaml.v2"
)

type LocalDesignDoc struct {
	Contents DesignDocContents
}

func NewLocalDesignDoc(fileName string) (*LocalDesignDoc, error) {
	if !strings.HasSuffix(fileName, ".yaml") {
		err := errors.Errorf("Required .yaml extension")
		log.WithError(err).WithField("file", fileName).Error("Wrong design doc extension. (.yaml required)")
		return nil, err
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithError(err).WithField("file", fileName).Error("Could not read file!")
		return nil, errors.Annotatef(err, "Could not read file: %s", fileName)
	}

	ldd := &LocalDesignDoc{}
	if err := yaml.Unmarshal(data, &ldd.Contents); err != nil {
		log.WithError(err).WithField("file", fileName).Error("Could not parse YAML file!")
		return nil, errors.Annotatef(err, "Could not parse YAML file: %s", data)
	}
	ldd.Contents.FullId = fileName

	return ldd, nil
}

func (d *LocalDesignDoc) Id() string {
	parts := strings.Split(d.Contents.FullId, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".yaml")
}

func (d *LocalDesignDoc) Md5() string {
	return DesignDocMd5(d.Contents)
}
