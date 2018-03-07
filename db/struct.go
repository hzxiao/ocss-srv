package db

import "github.com/hzxiao/goutil"

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Icon     string `bson:"icon" json:"icon"`
	Role     int    `bson:"role" json:"role"`
	Status   int    `bson:"status" json:"status"`
	Create   int64  `bson:"create" json:"create"`
	Update   int64  `bson:"update" json:"update"`
}

type Dept struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Major struct {
	ID     string `bson:"_id" json:"id"`
	DeptID string `bson:"deptId" json:"deptId"`
	Name   string `bson:"name" json:"name"`
}

type Student struct {
	ID         string     `bson:"_id" json:"id"` //学号
	Name       string     `bson:"name" json:"name"`
	Dept       goutil.Map `bson:"dept" json:"dept"`
	Major      goutil.Map `bson:"major" json:"major"`
	Clazz      string     `bson:"clazz" json:"clazz"`
	Sex        string     `bson:"sex" json:"sex"`
	Age        int        `bson:"age" json:"age"`
	Credit     string     `bson:"credit" json:"credit"`
	Email      string     `bson:"email" json:"email"`
	Phone      string     `bson:"phone" json:"phone"`
	SchoolYear string     `bson:"schoolYear" json:"schoolYear"`
	Create     int64      `bson:"create" json:"create"`
	Update     int64      `bson:"update" json:"update"`
	Status     int        `bson:"status" json:"status"`
}

type Teacher struct {
	ID         string     `bson:"_id" json:"id"` //工号
	Name       string     `bson:"name" json:"name"`
	Dept       goutil.Map `bson:"dept" json:"dept"`
	Sex        string     `bson:"sex" json:"sex"`
	Age        int        `bson:"age" json:"age"`
	Credit     string     `bson:"credit" json:"credit"`
	Email      string     `bson:"email" json:"email"`
	Phone      string     `bson:"phone" json:"phone"`
	SchoolYear string     `bson:"schoolYear" json:"schoolYear"`
	Title      string     `bson:"title" json:"title"` //职称
	Create     int64      `bson:"create" json:"create"`
	Update     int64      `bson:"update" json:"update"`
	Status     int        `bson:"status" json:"status"`
}

type Course struct {
	ID      string     `bson:"_id" json:"id"` //课程号
	Name    string     `bson:"name" json:"name"`
	Dept    goutil.Map `bson:"dept" json:"dept"`
	Credit  string     `bson:"credit" json:"credit"` //学分
	Period  string     `bson:"period" json:"period"` //学时
	Attr    string     `bson:"attr" json:"attr"`     //归属
	Nature  string     `bson:"nature" json:"nature"` //性质
	Campus  string     `bson:"campus" json:"campus"` //校区
	Desc    string     `bson:"desc" json:"desc"`     //描述
	Status  int        `bson:"status" json:"status"`
	Create  int64      `bson:"create" json:"create"`
	Update  int64      `bson:"update" json:"update"`
	Creator string     `bson:"creator" json:"creator"` //创建者
}

//授课
type TeachCourse struct {
	ID  string `bson:"_id" json:"id"`
	CID string `bson:"cid" json:"cid"` //course id
	TID string `bson:"tid" json:"tid"` //teacher id
	//eg:
	//[
	//	{
	//		"startWeek": 1,
	//		"endWeek": 8
	//	},
	//	{
	//		"startWeek": 12,
	//		"endWeek": 18
	//	}
	//]
	TakeWeeks []goutil.Map `bson:"takeWeeks" json:"takeWeeks"` //
	//eg:
	//[
	//	{
	//		"dayOfWeek": 1,
	//		"startTime": "10:00",
	//		"endTime": "12:00"
	//	},
	//	{
	//		"dayOfWeek": 2,
	//		"startTime": "10:00",
	//		"endTime": "12:00"
	//	}
	//]
	TakeTimes       []goutil.Map `bson:"takeTimes" json:"takeTimes"`
	StartSelectTime int64        `bson:"startSelectTime" json:"startSelectTime"`
	EndSelectTime   int64        `bson:"endSelectTime" json:"endSelectTime"`
	//eg: {"building": "XX楼", "classroom": "201"}
	Addr     goutil.Map `bson:"addr" json:"addr"`
	Capacity int        `bson:"capacity" json:"capacity"`
	Margin   int        `bson:"margin" json:"margin"`
	Status   int        `bson:"status" json:"status"`
	Create   int64      `bson:"create" json:"create"`
	Update   int64      `bson:"update" json:"update"`
}

//选课
type LearnCourse struct {
	ID            string  `bson:"_id" json:"id"`
	TCID          string  `bson:"tcid" json:"tcid"` //teach course id
	SID           string  `bson:"sid" json:"sid"`   //student id
	Grade         float64 `bson:"grade" json:"grade"`
	OrdinaryGrade float64 `bson:"ordinaryGrade" json:"ordinaryGrade"`
	ExamGrade     float64 `bson:"examGrade" json:"examGrade"`
	Create        int64   `bson:"create" json:"create"`
	Update        int64   `bson:"update" json:"update"`
}

type CourseResource struct {
	ID     string `bson:"_id" json:"id"`
	TCID   string `bson:"tcid" json:"tcid"` //teach course id
	TID    string `bson:"tid" json:"tid"`   //teacher id
	Desc   string `bson:"desc" json:"desc"` //描述
	File   *File  `bson:"file" json:"file"`
	Status int    `bson:"status" json:"status"`
	Create int64  `bson:"create" json:"create"`
	Update int64  `bson:"update" json:"update"`
}

type File struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Size int64  `bson:"size" json:"size"`
	Ext  string `bson:"ext" json:"ext"`
	Url  string `bson:"url" json:"url"`
}

type Comment struct {
	ID      string `bson:"_id" json:"id"`
	TCID    string `bson:"tcid" json:"tcid"` //teach course id
	UID     string `bson:"uid" json:"uid"`
	Role    int    `bson:"role" json:"role"`
	Content string `bson:"content" json:"content"`
	Create  int64  `bson:"create" json:"create"`
	Status  int    `bson:"status" json:"status"`
	//{
	//	"id": "xs",
	//	"uid": "123",
	//	"role": 1,
	//	"content": "xx",
	//	"to": "xxx",
	//	"create": "12345678"
	//}
	Children []goutil.Map `bson:"children" json:"children"` //子评论
}
