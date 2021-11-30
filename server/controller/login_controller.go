package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
)

type LoginReqBody struct {
	Key     string   `json:"key"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type LoginController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
}

type LoginResponse struct {
	ShuttingDown  bool `json:"login"`
	InProgressOps int  `json:"in_progress_operations"`
}

type LogoutResponse struct {
	ShuttingDown  bool `json:"logout"`
	InProgressOps int  `json:"in_progress_operations"`
}

func (d *LoginController) Login(ctx iris.Context) {

}

func (d *LoginController) Logout(ctx iris.Context) {

}
