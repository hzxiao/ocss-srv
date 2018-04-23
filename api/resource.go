package api

import (
	"bytes"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/config"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func UploadFile(ctx context.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	log.Printf("[UploadFile] add file(%v) fileHeader(%v)", file, fileHeader)
	f := &db.File{}
	_, f.Name = filepath.Split(fileHeader.Filename)
	f.Ext = strings.ToLower(filepath.Ext(f.Name))
	f.ID = tools.GenerateUniqueId() + f.Ext
	f.Size = fileHeader.Size
	f.Url = "/files/" + f.ID
	err = tools.SaveFile(config.GetString("file.location"), f.ID, file)
	if err != nil {
		log.Errorf("[UploadFile] save file(%v) error(%v)", goutil.Struct2Json(f), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	err = db.AddFile(f)
	if err != nil {
		log.Errorf("[UploadFile] add file(%v) error(%v)", goutil.Struct2Json(f), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"file": f,
	})
}

func CallUploadFile(filename string) (goutil.Map, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}
	w.Close()
	result, err := tools.HttpPost(fmt.Sprintf("http://%v/files", SrvAddr), "",
		w.FormDataContentType(), &b)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func GetFile(ctx context.Context) {
	id := ctx.Params().Get("id")
	dl := ctx.FormValue("dl")

	file, err := db.LoadFile(id)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}
		log.Errorf("[GetFile] get file(%v) error(%v)", id, err)
		WriteResultWithSrvErr(ctx, err)
	}
	filename := config.GetString("file.location") + string(filepath.Separator) + file.ID
	if dl == "1" {
		err = ctx.SendFile(filename, file.Name)
	} else {
		err = ctx.ServeFile(filename, false)
	}

	if err != nil {
		log.Errorf("[GetFile] send file(%v) dl(%v) error(%v)", id, dl, err)
	}
}

func AddCourseResource(ctx context.Context) {
	var r db.CourseResource
	err := ctx.ReadJSON(&r)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	r.TID = ctx.Values().GetString("uid")
	err = db.AddResource(&r)
	if err != nil {
		log.Errorf("[AddCourseResource] add resource(%v) error(%v)", goutil.Struct2Json(r), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"resource": r,
	})
}

func DelCourseResource(ctx context.Context) {
	var ids []string
	err := ctx.ReadJSON(&ids)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	var tid string
	role, err := ctx.Values().GetFloat64("role")
	if err != nil {
		log.Errorf("[DelCourseResource] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	if int(role) != db.RoleAdmin {
		tid = ctx.Values().GetString("uid")
	}
	err = db.DelCourseResource(tid, "", ids)
	if err != nil {
		log.Errorf("[DelCourseResource] delete ids(%v) error(%v)", ids, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"ids": ids,
	})
}

func GetCourseResource(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "tcid", Type: "string"},
		{Key: "tid", Type: "string"},
		{Key: "status", Type: "int", DefaultValue: strconv.Itoa(db.StatusNormal)},
		{Key: "page", Type: "int", DefaultValue: "1"},
		{Key: "pageSize", Type: "int", DefaultValue: "20"},
		{Key: "sort", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	cond := TakeByKeys(argMap, "tcid", "tid", "status")
	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}

	resourceList, total, err := db.ListCourseResource(cond, sort, skip, limit)
	if err != nil {
		log.Errorf("[GetCourseResource] get by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}


	var resList []goutil.Map

	for i := range resourceList {
		res := goutil.Struct2Map(resourceList[i])
		tid := resourceList[i].TID
		tcid := resourceList[i].TCID
		teacher,err := db.LoadTeacher(tid)
		if err != nil {
			log.Errorf("[GetCourseResource] get by(%v) error(%v)", teacher, err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		res.Set("teacherName", "")
		res.Set("courseName", "")
		if teacher != nil {
			res.Set("teacherName", teacher.Name)
		}
		teachCourse,err := db.LoadTeachCourse(tcid)
		if err != nil {
			log.Errorf("[GetCourseResource] get by(%v) error(%v)", teachCourse, err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		if teachCourse != nil {
			course,err := db.LoadCourse(teachCourse.CID)
			if err != nil {
				log.Errorf("[GetCourseResource] get by(%v) error(%v)", course, err)
				WriteResultWithSrvErr(ctx, err)
				return
			}
			if course != nil {
				res.Set("courseName", course.Name)
			}
		}
		resList = append(resList,res)

	}
	log.Printf("[GetCourseResource] by argMap(%)  get resList(%v) ", argMap, resList)
	
	WriteResultSuccess(ctx, goutil.Map{
		"resourceList": resList,
		"total":        total,
	})
}
