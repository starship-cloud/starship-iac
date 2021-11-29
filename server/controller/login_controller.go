package controllers

import (
	"encoding/json"
	"fmt"
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
	ShuttingDown  bool `json:"shutting_down"`
	InProgressOps int  `json:"in_progress_operations"`
}

func (d *LoginController) Login(ctx iris.Context) {
	status := d.Drainer.GetStatus()
	data, err := json.MarshalIndent(&UsersResponse{
		ShuttingDown:  status.ShuttingDown,
		InProgressOps: status.InProgressOps,
	}, "", "  ")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		d.Logger.Info(fmt.Sprintf("Error creating user json response: %s", err))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(data)
}

func (d *LoginController) Logout(ctx iris.Context) {
	status := d.Drainer.GetStatus()
	data, err := json.MarshalIndent(&UsersResponse{
		ShuttingDown:  status.ShuttingDown,
		InProgressOps: status.InProgressOps,
	}, "", "  ")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		d.Logger.Info(fmt.Sprintf("Error creating user json response: %s", err))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(data)
}
