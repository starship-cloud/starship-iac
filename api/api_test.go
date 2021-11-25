package api

import (
	"github.com/kataras/iris/v12"
	"strconv"
	"testing"
)

func Test_Token(t *testing.T) {
	app := Init()
	app.Get("/", test_login)
	app.Get("/protected", test_check)
	app.Run(iris.Addr(":8888"))
}

func test_login(ctx iris.Context) {
	tokenString, _ := createToken("abcd")
	ctx.HTML(`Token: ` + tokenString + `<br/><br/><a href="/protected?token=` + tokenString + `">/protected?token=` + tokenString + `</a>`)
}

func test_check(ctx iris.Context) {
	tokenString, isUseful := checkToken("abcd", ctx)
	ctx.HTML(`Token: ` + tokenString + `,useful:` + strconv.FormatBool(isUseful) + `<br/><br/><a href="/protected?token=` + tokenString + `">/protected?token=` + tokenString + `</a>`)
}
