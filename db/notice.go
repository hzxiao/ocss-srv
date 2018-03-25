package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

var (
	AdminNoticeChan   chan goutil.Map
	StudentNoticeChan chan goutil.Map
	TeacherNoticeChan chan goutil.Map
)

func init() {
	AdminNoticeChan = make(chan goutil.Map)
	StudentNoticeChan = make(chan goutil.Map)
	TeacherNoticeChan = make(chan goutil.Map)

	go receiveNotice()
}

func SendNotice(role int, msg goutil.Map) {
	switch role {
	case RoleAdmin:
		AdminNoticeChan <- msg
	case RoleTeacher:
		TeacherNoticeChan <- msg
	case RoleStudent:
		StudentNoticeChan <- msg
	}
}

func receiveNotice() {
	for {
		select {
		case msg := <-AdminNoticeChan:
			_ = msg
		case msg := <-TeacherNoticeChan:
			_ = msg

		case msg := <-StudentNoticeChan:
			_ = msg
		}
	}
}

func AddNotice(ns ...*Notice) error {
	if len(ns) == 0 {
		return nil
	}

	var docs []interface{}
	for i := range ns {
		if ns[i].UID == "" {
			return errors.New("the %vth notice's uid is empty")
		}
		ns[i].ID = tools.GenerateUniqueId()
		ns[i].Status = NoticeStateUnRead
		ns[i].Create = tools.NowMillisecond()
		ns[i].Update = ns[i].Create

		docs = append(docs, ns[i])
	}
	return C(CollectionNotice).Insert(docs...)
}

func ListNotice(cond goutil.Map, sort []string, skip, limit int) ([]*Notice, int, error) {
	var noticeList []*Notice
	total, err := list(CollectionNotice, tools.ToBsonMap(cond), nil, sort, skip, limit, &noticeList)
	if err != nil {
		return nil, 0, err
	}
	return noticeList, total, nil
}

func UpdateNotice(notice *Notice) error {
	if notice == nil {
		return errors.New("notice is nil")
	}

	if notice.ID == "" {
		return errors.New("id is empty")
	}

	args := bson.M{}
	if notice.Status > 0 {
		args["status"] = notice.Status
	}

	args["update"] = tools.NowMillisecond()

	return C(CollectionNotice).UpdateId(notice.ID, bson.M{"$set": args})
}

func CountNoticeDiffStatus(uid string) (goutil.Map, error) {
	pipe := []bson.M{
		{"$match": bson.M{"uid": uid}},
		{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
	}
	var list []goutil.Map
	err := C(CollectionNotice).Pipe(pipe).All(&list)
	if err != nil {
		return nil, err
	}

	res := goutil.Map{}
	for _, vm := range list {
		switch int(vm.GetInt64("_id")) {
		case NoticeStatusRead:
			res.Set("read", vm.GetInt64("count"))
		case NoticeStateUnRead:
			res.Set("unread", vm.GetInt64("count"))
		case NoticeStatusDeleted:
			res.Set("deleted", vm.GetInt64("count"))
		}
	}
	return res, nil
}