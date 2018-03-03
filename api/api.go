package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hzxiao/goutil"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/juju/errors"
	"github.com/kataras/iris"
)

var SrvAddr string

func RegisterHandle(app *iris.Application) {

	//common
	app.Post("/login", Login)
	//user
	userRouter := app.Party("/users")
	UseJwt(userRouter)

	userRouter.Post("/", AddUser)
	userRouter.Put("/{username:string}", UpdateUser)
	userRouter.Get("/{username:string}", GetUser)
}

func UseJwt(partys ...iris.Party) {
	JwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("ocss"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	for i := range partys {
		partys[i].Use(JwtMiddleware.Serve)
		partys[i].Use(func(ctx iris.Context) {
			userToken := JwtMiddleware.Get(ctx)
			if claims, ok := userToken.Claims.(jwt.MapClaims); ok && userToken.Valid {
				ctx.Values().Set("uid", claims["uid"])
				ctx.Values().Set("uid", claims["role"])
				ctx.Next()
			} else {
				ctx.StatusCode(iris.StatusUnauthorized)
			}
		})
	}
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
