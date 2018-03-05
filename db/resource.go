package db

import (
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

func AddFile(file *File) error {
	if file == nil {
		return errors.New("file is nil")
	}
	if file.ID == "" {
		return errors.New("id is empty")
	}

	return C(CollectionFile).Insert(file)
}

func LoadFile(id string) (*File, error) {
	var file File
	err := one(CollectionFile, bson.M{"_id": id}, nil, &file)
	return &file, err
}
