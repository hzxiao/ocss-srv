package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
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
			//finder admin
			adminList, err := ListUserByRole(RoleAdmin)
			if err != nil {
				log.Error(err)
				break
			}
			var ns []*Notice
			for i := range adminList {
				ns = append(ns, &Notice{
					UID:     adminList[i].ID,
					Content: msg.GetString("content"),
					Title:   msg.GetString("title"),
				})
			}
			err = AddNotice(ns...)
			if err != nil {
				log.Error(err)
				break
			}
		case msg := <-TeacherNoticeChan:
			tea := msg.GetString("tid")

			var ns []*Notice
			ns = append(ns, &Notice{
				UID:     tea,
				Content: msg.GetString("content"),
				Title:   msg.GetString("title"),
			})
			err := AddNotice(ns...)
			if err != nil {
				log.Error(err)
				break
			}

		case msg := <-StudentNoticeChan:
			stuList := msg.GetStringArray("sid")
			var ns []*Notice
			for i := range stuList {
				ns = append(ns, &Notice{
					UID:     stuList[i],
					Content: msg.GetString("content"),
					Title:   msg.GetString("title"),
				})
			}
			log.Printf("receiveNotice: stu(%v)", stuList)

			err := AddNotice(ns...)
			if err != nil {
				log.Error(err)
				break
			}
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
