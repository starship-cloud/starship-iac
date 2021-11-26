package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
)

type RequestBody struct {
	Key     string   `json:"key"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

// StatusController handles the status of Atlantis.
type StatusController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
}

type StatusResponse struct {
	ShuttingDown  bool `json:"shutting_down"`
	InProgressOps int  `json:"in_progress_operations"`
}

func (d *StatusController) Status(ctx iris.Context) int {
	status := d.Drainer.GetStatus()
	data, err := json.MarshalIndent(&StatusResponse{
		ShuttingDown:  status.ShuttingDown,
		InProgressOps: status.InProgressOps,
	}, "", "  ")
	if err != nil {
		d.Logger.Info(fmt.Sprintf("Error creating status json response: %s", err))
		return iris.StatusInternalServerError
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(data)

	return iris.StatusOK
}
