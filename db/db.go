package db

import (
	"gopkg.in/mgo.v2"
)

var C = func(name string) *mgo.Collection {
	panic("the database collection handle function is not initial")
}

func InitDB(url string, dbName string) error {
	sess, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	C = sess.DB(dbName).C
	return nil
}

//db collection name
const (
	CollectionUser = "ocss_user"
)
