package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func AddComment(c *Comment) error {
	if c == nil {
		return errors.New("c is nil")
	}
	if c.TCID == "" {
		return errors.New("tcid is nil")
	}
	if c.UID == "" {
		return errors.New("uid is empty")
	}
	c.ID = tools.GenerateUniqueId()
	c.Status = StatusNormal
	c.Create = tools.NowMillisecond()
	return C(CollectionComment).Insert(c)
}

func ListComment(cond goutil.Map, sort []string, skip, limit int) ([]*Comment, int, error) {
	var commentList []*Comment
	total, err := list(CollectionComment, tools.ToBsonMap(cond), nil, sort, skip, limit, &commentList)
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

func UpdateChildComment(id string, method string, child goutil.Map) (*Comment, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	if len(child) == 0 {
		return nil, errors.New("child is empty")
	}
	updater := bson.M{}
	switch method {
	case "add":
		child.Set("id", tools.GenerateUniqueId())
		child.Set("create", tools.NowMillisecond())
		updater["$push"] = bson.M{"children": child}
	case "del":
		if child.GetString("id") == "" {
			return nil, errors.New("del child comment expect id, but id is empty")
		}
		updater["$pull"] = bson.M{"children": bson.M{"id": child.GetString("id")}}
	default:
		return nil, errors.New("unknown method")
	}

	updater["$set"] = bson.M{"update": tools.NowMillisecond()}
	var c Comment
	_, err := C(CollectionComment).FindId(id).Apply(mgo.Change{
		Update:    updater,
		ReturnNew: true,
	}, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func DelComment(ids []string) error {
	_, err := C(CollectionComment).UpdateAll(bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$set": bson.M{"status": StatusDeleted}})

	return err
}
