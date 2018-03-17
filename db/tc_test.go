package db

import (
	"testing"
	"github.com/hzxiao/goutil/assert"
	"gopkg.in/mgo.v2/bson"
	"github.com/hzxiao/goutil"
)

func TestAddTeachCourse(t *testing.T) {
	C(CollectionTeachCourse).RemoveAll(nil)
	teachCourse:=&TeachCourse{
		TID:"t1",
		CID:"c1",
	}
	err:=AddTeachCourse(teachCourse)
	assert.NoError(t,err)

	teachCourse=&TeachCourse{
		TID:"t1",
		CID:"c2",
	}
	err=AddTeachCourse(teachCourse)
	assert.NoError(t,err)
	teachCourse=&TeachCourse{
		TID:"t1",
		CID:"c3",
	}
	err=AddTeachCourse(teachCourse)
	assert.NoError(t,err)

	//列出可选课程
	listTCourse,total,err:=ListTeachCourse(1,nil,0,0)
	assert.NoError(t,err)
	assert.NotNil(t,listTCourse)
	assert.Equal(t,3,total)
	//t.Log(total)
	//学生选课
	var tcids []string
	for _,v:= range listTCourse  {
		tcids = append(tcids,v.ID)
	}
	err=StuSelectCourse(tcids,"s1")
	assert.NoError(t,err)
	var  tcs []*TeachCourse
	C(CollectionTeachCourse).Find(bson.M{"_id":bson.M{"$in":tcids}, "status":TeachCourseStatusSelectable}).All(&tcs)
	//t.Log(goutil.Struct2Json(tcs))

	listSTCourse,total,err:=ListStuLearnCourse(1,"s1",nil,0,0)
	assert.NoError(t,err)
	assert.NotNil(t,listSTCourse)
	assert.Equal(t,3,total)


	err=StuCancelCourse(tcids,"s1")
	assert.NoError(t,err)

	listSTCourse,total,err=ListStuLearnCourse(1,"s1",nil,0,0)
	assert.NoError(t,err)
	assert.Nil(t,listSTCourse)
	assert.Equal(t,0,total)


	tcids1:=[]string{tcids[0]}
	err=StuSelectCourse(tcids1,"s1")
	assert.NoError(t,err)
	tcids2:=[]string{tcids[1]}
	err=StuSelectCourse(tcids2,"s2")
	assert.NoError(t,err)
	err=StuSelectCourse(tcids2,"s3")
	assert.NoError(t,err)
	sids:=[]string{"s1","s2"}
	listTTCourse,total,err:=TeaListStuOfCourse("t1",sids,1,nil,0,0)
	assert.NoError(t,err)
	assert.NotNil(t,listTTCourse)
	assert.Equal(t,2,total)

	//通过选课
	sinfo:=[]goutil.Map{
		goutil.Map{
			"sid":"s2",
			"cstatus":SCourseStatusSelected,
		},
		goutil.Map{
			"sid":"s3",
			"cstatus":SCourseStatusSelected,
		},
	}
	err = TeaSettingGrade(tcids[1],sinfo)
	assert.NoError(t,err)
	//课程进入学习状态
	tcourse:=&TeachCourse{Status:TeachCourseStatusLearning}
	err = UpdateTeachCourseByIDs(tcids2,tcourse)
	assert.NoError(t,err)


	//评分
	info:=[]goutil.Map{
		goutil.Map{
			"sid":"s2",
			"grade":99,
		},
		goutil.Map{
			"sid":"s3",
			"grade":88,
		},
	}


	err = TeaSettingGrade(tcids[1],info)
	assert.NoError(t,err)
}
