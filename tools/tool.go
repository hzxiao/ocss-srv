package tools

import "gopkg.in/mgo.v2/bson"

func GenerateUniqueId() string {
	return bson.NewObjectId().Hex()
}