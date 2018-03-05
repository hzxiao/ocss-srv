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
	"strings"
)

func UploadFile(ctx context.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	f := &db.File{}
	_, f.Name = filepath.Split(fileHeader.Filename)
	f.Ext = strings.ToLower(filepath.Ext(f.Name))
	f.ID = tools.GenerateUniqueId() + f.Ext
	f.Size = fileHeader.Size
	f.Url = "/files/" + f.ID
	err = tools.SaveFile(config.GetString("file.location"), f.ID, file)
	if err != nil {
		log.Printf("[UploadFile] save file(%v) error(%v)", goutil.Struct2Json(f), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	err = db.AddFile(f)
	if err != nil {
		log.Printf("[UploadFile] add file(%v) error(%v)", goutil.Struct2Json(f), err)
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
		log.Printf("[GetFile] get file(%v) error(%v)", id, err)
		WriteResultWithSrvErr(ctx, err)
	}
	filename := config.GetString("file.location") + string(filepath.Separator) + file.ID
	if dl == "1" {
		err = ctx.SendFile(filename, file.Name)
	} else {
		err = ctx.ServeFile(filename, false)
	}

	if err != nil {
		log.Printf("[GetFile] send file(%v) dl(%v) error(%v)", id, dl, err)
	}
}
