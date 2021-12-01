package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
)

type requestBody struct {
	Key     string   `json:"key"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type OperationController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
}

type ApplyResponse struct {
	ShuttingDown  bool `json:"shutting_down"`
	InProgressOps int  `json:"in_progress_operations"`
}

func (d *OperationController) Status(ctx iris.Context) {

	return
}
