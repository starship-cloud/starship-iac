package api

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/cmd"
	"github.com/starship-cloud/starship-iac/file"
	"github.com/starship-cloud/starship-iac/taskpool"
	"github.com/starship-cloud/starship-iac/utils"
	"log"
)

//var app *iris.Application
func Init() *iris.Application{
	app := iris.New()
	app.Post("/apply", apply)
	app.Post("/cancel", cancel)
	return app
}

type RequestBody struct {
	Key     string   `json:"key"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
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
