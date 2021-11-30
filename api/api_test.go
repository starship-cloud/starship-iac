package api

import (
	"github.com/kataras/iris/v12"
	"testing"
)

func Test_Token(t *testing.T) {
	app := Init()
	app.Get("/login", test_login)
	app.Get("/protected", test_check)
	app.UseGlobal(checkToken)
	app.Run(iris.Addr(":8888"))
}

func test_login(ctx iris.Context) {
	tokenString, _ := createToken("abcd")
	ctx.JSON("token:" + tokenString)
}

func test_check(ctx iris.Context) {
	ctx.HTML("success")
}
