package controllers

import (
	"github.com/casbin/casbin/v2"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	service "github.com/starship-cloud/starship-iac/server/services"
)

type PermissionController struct {
	Logger   logging.SimpleLogging
	Drainer  *events.Drainer
	Enforcer *casbin.Enforcer
}

type PermissionResp struct {
	StatusCode  uint
	Description string
}

func (pc *PermissionController) AddPermission(ctx iris.Context) {
	var permission models.Permission
	ctx.ReadJSON(&permission)
	_, err := service.AddPermission(permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "authorize failed.",
		})
	}
}
