package db

import (
	"github.com/hzxiao/ocss-srv/config"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"github.com/hzxiao/ocss-srv/tools"
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

	//init data
	err = InitDept(config.GetString("data.deptPath"))
	if err != nil {
		panic(err)
	}
	err = InitMajor(config.GetString("data.majorPath"))
	if err != nil {
		panic(err)
	}

	go PingLoop(sess, url, dbName)

	return nil
}

//db collection name
const (
	CollectionUser        = "user"
	CollectionDept        = "dept"
	CollectionMajor       = "major"
	CollectionStudent     = "student"
	CollectionTeacher     = "teacher"
	CollectionCourse      = "course"
	CollectionFile        = "file"
	CollectionResource    = "resource"
	CollectionComment     = "comment"
	CollectionNotice      = "notice"
	CollectionTeachCourse = "teach_course"
)

func PingLoop(sess *mgo.Session, url, dbName string) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		<-ticker.C
		err := Ping(sess)
		if err == nil {
			log.Printf("ping to mongo success by url(%v) db(%v)\n", url, sess.DB(dbName).Name)
			continue
		}
		//handle err
		for {
			sess, err = mgo.Dial(url)
			if err != nil {
				log.Printf("try to dial mongo by url(%v) fail. \n", url)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("reconnect to mongo success by url(%v)\n", url)
			C = sess.DB(dbName).C
			break
		}
	}
}

func Ping(sess *mgo.Session) (err error) {
	errClosed := errors.New("Closed explicitly")
	defer func() {
		if pe := recover(); pe != nil {
			if sess != nil {
				sess.Clone()
				err = errClosed
			}
		}
	}()

	err = sess.Ping()
	if err == nil {
		return nil
	}
	if err.Error() == "Closed explicitly" || err.Error() == "EOF" {
		sess.Clone()
		return errClosed
	}
	return err
}

func PrepareData() error {
	fileDir := "../data/"

	//var users []*User
	//err := tools.UnmarshalJsonFile(fileDir+"users.json", &users)
	//if err != nil {
	//	return err
	//}
	//
	//for i := range users {
	//	err = AddUser(users[i])
	//	if err != nil {
	//		return err
	//	}
	//}

	var courses []*Course
	err := tools.UnmarshalJsonFile(fileDir+"course.json", &courses)
	if err != nil {
		return err
	}

	for i := range courses {
		_, err = AddCourse(courses[i])
		if err != nil {
			return err
		}
	}

	var teachers []*Teacher
	err = tools.UnmarshalJsonFile(fileDir+"teachers.json", &teachers)
	if err != nil {
		return err
	}

	for i := range teachers {
		err = AddTeacher(teachers[i])
		if err != nil {
			return err
		}
	}

	var tcs []*TeachCourse
	err = tools.UnmarshalJsonFile(fileDir+"tc.json", &tcs)
	if err != nil {
		return err
	}

	for i := range tcs {
		err = AddTeachCourse(tcs[i])
		if err != nil {
			return err
		}
	}

	return nil
}