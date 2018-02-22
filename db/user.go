package db

import (
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

func AddUser(user *User) error {
	if user == nil {
		return errors.New("user can not be nil")
	}
	user.ID = bson.NewObjectId().Hex()
	return C(CollectionUser).Insert(user)
}
