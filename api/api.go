package api

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/cmd"
	"github.com/starship-cloud/starship-iac/file"
	"github.com/starship-cloud/starship-iac/taskpool"
	"github.com/starship-cloud/starship-iac/utils"
	"log"
	"net/http"
	"time"
)

type RequestBody struct {
	Key     string   `json:"key"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

var jwtMiddleware = jwt.New(jwt.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.RootSecret), nil
	},
	Expiration:    true,
	Extractor:     jwt.FromParameter("token"),
	SigningMethod: jwt.SigningMethodHS256,
})

func Init() *iris.Application {
	app := iris.New()
	app.UseGlobal(checkToken)
	app.Post("/apply", apply)
	app.Post("/cancel", cancel)
	return app
}

func createToken(userId string) (string, error) {
	now := time.Now()
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iat":    now.Unix(),
		"exp":    now.Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(utils.RootSecret))
}

func checkToken(ctx iris.Context) {
	path := ctx.Path()
	if path == "/login" || path == "/healthz" || path == "/status" {
		ctx.Next()
		return
	}

	if err := jwtMiddleware.CheckJWT(ctx); err != nil {
		jwtMiddleware.Config.ErrorHandler(ctx, err)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "Wrong token")
		return
	}

	token := ctx.Values().Get("jwt").(*jwt.Token)
	tokenInfo := token.Claims.(jwt.MapClaims)
	userId := ctx.URLParam("id")
	if userId != tokenInfo["userId"] {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "User does not have permission.")
		return
	}
	if time.Now().Unix() > int64(tokenInfo["exp"].(float64)) {
		//token timeout
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Values().Set("msg", "Token timeout.")
		return
	} else {
		//update token
		if int64(tokenInfo["exp"].(float64))-time.Now().Unix() < 30 {
			newToken, _ := createToken(userId)
			ctx.Header("token", newToken)
			ctx.Next()
			return
		}
	}
	oldToken, _ := token.SignedString([]byte(utils.RootSecret))
	ctx.Header("token", oldToken)
	ctx.Next()
}

func apply(ctx iris.Context) {
	var rb RequestBody
	ctx.ReadJSON(&rb)
	task := taskpool.Task{}
	task.Do = func() {
		task.Command = cmd.Exec(rb.Command, rb.Params)
		file.WriteFileByCmd("test.log", utils.RootCmdLogPath, task.Command)
	}
	task.Stop = func() {
		if err := task.Command.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}
	taskpool.Run(rb.Key, task)
	ctx.JSON("apply finish")
}

func cancel(ctx iris.Context) {
	var rb RequestBody
	ctx.ReadJSON(&rb)
	taskpool.Cancel(rb.Key)
	ctx.JSON("cancel finish")
}

func getLog(ctx iris.Context) {
	lines, total := cmd.ReadLog(utils.RootCmdLogPath+"test.log", 6)
	ctx.JSON(total)
	ctx.JSON(lines)
}
