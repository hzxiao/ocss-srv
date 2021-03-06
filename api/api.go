package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hzxiao/goutil"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"strings"
)

var SrvAddr string

func RegisterHandle(app *iris.Application) {

	//common
	app.Post("/login", Login)
	app.Post("/files", UploadFile)
	app.Get("/files/{id:string}", GetFile)
	app.Post("/count", Count)
	//users
	userRouter := app.Party("/users")
	UseJwt(userRouter)
	userRouter.Post("/", AddUser)
	userRouter.Put("/{username:string}", UpdateUser)
	userRouter.Get("/{username:string}", GetUser)
	userRouter.Post("/password", UpdateUserPassword)

	//students
	stuRouter := app.Party("/students")
	UseJwt(stuRouter)
	stuRouter.Post("/add", AddStudent)
	stuRouter.Put("/{id:string}", UpdateStudent)
	stuRouter.Delete("/delete/", DeleteStudent)
	stuRouter.Get("/list", GetStudents)
	stuRouter.Get("/{id:string}", GetStudent)
	stuRouter.Get("/checkImport", CheckImportStudent)
	stuRouter.Get("/import", ImportStudent)

	//teachers
	teacherRouter := app.Party("/teachers")
	UseJwt(teacherRouter)
	teacherRouter.Post("/add", AddTeacher)
	teacherRouter.Put("/{id:string}", UpdateTeacher)
	teacherRouter.Delete("/delete", DeleteTeacher)
	teacherRouter.Get("/list", GetTeachers)
	teacherRouter.Get("/{id:string}", GetTeacher)

	//depts
	app.Get("/depts", GetAllDept)

	//majors
	app.Get("/majors", GetAllMajor)

	//resource
	resRouter := app.Party("/resources")
	UseJwt(resRouter)
	resRouter.Post("/add", AddCourseResource)
	resRouter.Delete("/delete", DelCourseResource)
	resRouter.Get("/list", GetCourseResource)

	//comments
	commentRouter := app.Party("/comments")
	UseJwt(commentRouter)
	commentRouter.Post("/add", AddComment)
	commentRouter.Delete("/delete/{id:string}", DelComment)
	commentRouter.Get("/list", ListComment)
	commentRouter.Post("/{id:string}/children", AddChildComment)
	commentRouter.Delete("/delete/{id:string}/children/{childId:string}", DelChildComment)

	//notices
	noticeRouter := app.Party("/notices")
	UseJwt(noticeRouter)
	noticeRouter.Put("/{id:string}", UpdateNotice)
	noticeRouter.Get("/list", ListNotice)
	noticeRouter.Get("/count", CountNotice)

	//course
	courseRouter := app.Party("/courses")
	UseJwt(courseRouter)
	courseRouter.Post("/add", AddCourse)
	courseRouter.Put("/{id:string}", UpdateCourse)
	courseRouter.Delete("/delete/", DeleteCourse)
	courseRouter.Get("/list", GetCourses)
	courseRouter.Get("/{id:string}", GetCourse)

	tcRouter := app.Party("/tc")
	UseJwt(tcRouter)
	tcRouter.Post("/add", AddTeachCourse)
	tcRouter.Put("/", UpdateTeachCourse)
	tcRouter.Put("/all", UpdateTeachCourses)
	tcRouter.Get("/list", ListTeachCourse)
	tcRouter.Get("/stu/grade/list", ListTeachCourseForStu)
	tcRouter.Delete("/", DeleteTeachCourse)
	tcRouter.Put("/updateGrade", SetGradeForTc)
	tcRouter.Get("/stu/courses", ListStudentCourse)
	tcRouter.Post("/stu/selectCourse", StuSelectCourse)
	tcRouter.Get("/{id:string}", GetTeachCourse)
	tcRouter.Get("/stu/list/{id:string}", ListStudentOfCourse)
	tcRouter.Post("/stu/update", UpdateStudentForTc)
}

func UseJwt(partys ...router.Party) {
	JwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("ocss"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	for i := range partys {
		partys[i].Use(JwtMiddleware.Serve)
		partys[i].Use(func(ctx context.Context) {
			userToken := JwtMiddleware.Get(ctx)
			if claims, ok := userToken.Claims.(jwt.MapClaims); ok && userToken.Valid {
				ctx.Values().Set("uid", claims["uid"])
				ctx.Values().Set("role", claims["role"])
				ctx.Next()
			} else {
				ctx.StatusCode(iris.StatusUnauthorized)
			}
		})
	}
}

func CheckPermission(ctx context.Context) {

}

func NewToken(uid string, role int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"role": role,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte("ocss"))

	return tokenString
}

func handleACallResult(result goutil.Map) (goutil.Map, error) {
	if result.GetInt64("code") == 0 {
		return result.GetMap("data"), nil
	}

	return result, errors.New(result.GetString("msg") + " " + result.GetString("err"))
}

func appendArgs(url string, argMap goutil.Map) string {
	var args []string
	for k := range argMap {
		args = append(args, k+"="+argMap.GetString(k))
	}
	if len(args) == 0 {
		return url
	}
	if !strings.HasSuffix(url, "?") {
		url = url + "?"
	}
	url = url + strings.Join(args, "&")
	return url
}
