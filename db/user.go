package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//AddUser add a new user by username, if username already exists, will return err
func AddUser(user *User) error {
	if user == nil {
		return errors.New("user can not be nil")
	}
	if user.Username == "" {
		return errors.New("user.username is empty")
	}

	initUser(user)

	changeInfo, err := C(CollectionUser).FindId(user.Username).Apply(mgo.Change{
		Update:    bson.M{"$setOnInsert": user},
		Upsert:    true,
		ReturnNew: true,
	}, user)
	if err != nil {
		log.WithFields(log.Fields{
			"user": goutil.Struct2Json(user),
			"err":  err,
		}).Error("AddUser find and modify err")
		return err
	}
	if changeInfo.UpsertedId == nil {
		return errors.New("username already exists")
	}
	return nil
}

func initUser(user *User) {
	if user.Password == "" {
		user.Password = DefaultPassword
	}
	if user.Status <= 0 || user.Status > 3 {
		user.Status = UserStatsNormal
	}
	if user.Role <= 0 || user.Role > 3 {
		user.Role = RoleStudent
	}

	user.ID = user.Username
	user.Create = tools.NowMillisecond()
	user.Update = tools.NowMillisecond()
}

//UpdateUser update user info by username
func UpdateUser(user *User) error {
	if user == nil {
		return errors.New("user can not be nil")
	}
	if user.Username == "" {
		return errors.New("user.username is empty")
	}

	args := bson.M{}
	if user.Password != "" {
		args["password"] = user.Password
	}
	if user.Role > 0 {
		args["role"] = user.Role
	}
	if user.Status > 0 {
		args["status"] = user.Status
	}
	if user.Icon != "" {
		args["icon"] = user.Icon
	}
	if len(args) == 0 {
		return nil
	}
	args["update"] = tools.NowMillisecond()
	_, err := C(CollectionUser).FindId(user.Username).Apply(mgo.Change{
		Update:    bson.M{"$set": args},
		ReturnNew: true,
	}, user)
	if err != nil {
		log.WithFields(log.Fields{
			"username": user.Username,
			"set":      goutil.Struct2Json(args),
			"err":      err,
		}).Error("UpdateUser update user info err")
		return err
	}
	return nil
}

//VerifyUser verify user by username and password,
//if the user is not found, return user and err are nil
func VerifyUser(username, password string) (*User, error) {
	var user User
	err := C(CollectionUser).Find(bson.M{"_id": username, "password": password}).One(&user)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		log.WithFields(log.Fields{
			"username": username,
			"password": password,
			"err":      err,
		}).Errorf("VerifyUser find user by username and password")
	}
	user.Password = ""
	return &user, nil
}

//FindUserByUsername find user by username
func FindUserByUsername(username string) (*User, error) {
	var user User
	err := C(CollectionUser).FindId(username).One(&user)
	return &user, err
}

func UpdateUserByIDs(ids []string, update goutil.Map) error {
	if len(ids) == 0 {
		return nil
	}
	if len(update) == 0 {
		return errors.New("update is empty")
	}

	_, err := C(CollectionUser).UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": tools.ToBsonMap(update)})
	return err
}

