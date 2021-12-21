package controllers

import (
	"github.com/casbin/casbin/v2"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	service "github.com/starship-cloud/starship-iac/server/services"
	"github.com/starship-cloud/starship-iac/utils"
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

func (pc *PermissionController) AddUserToRole(ctx iris.Context) {
	var role models.RoleForUser
	ctx.ReadJSON(&role)
	if role.RoleName != utils.Admin || role.RoleName != utils.ProjectCreator || role.RoleName != utils.SecretManager {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong role.",
		})
	}
	_, err := service.AddRoleForUser(&role, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add role failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add role success.",
		})
	}
}

func (pc *PermissionController) RemoveUserFromRole(ctx iris.Context) {
	var role models.RoleForUser
	ctx.ReadJSON(&role)
	if role.RoleName != utils.Admin || role.RoleName != utils.ProjectCreator || role.RoleName != utils.SecretManager {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong role.",
		})
	}
	_, err := service.DeleteRoleForUser(&role, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove role failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove role success.",
		})
	}
}

func (pc *PermissionController) AddUser(ctx iris.Context) {
	var permission models.ProjectPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.ReadOnly || permission.Permission != utils.Config {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong project permission.",
		})
	}
	_, err := service.AddProjectPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add permission success.",
		})
	}
}

func (pc *PermissionController) RemoveUser(ctx iris.Context) {
	var permission models.ProjectPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.ReadOnly || permission.Permission != utils.Config {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong project permission.",
		})
	}
	_, err := service.DeleteProjectPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove permission success.",
		})
	}
}

func (pc *PermissionController) AddGroup(ctx iris.Context) {
	var permission models.ProjectPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.ReadOnly || permission.Permission != utils.Config {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong project permission.",
		})
	}
	_, err := service.AddProjectPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add permission success.",
		})
	}
}

func (pc *PermissionController) RemoveGroup(ctx iris.Context) {
	var permission models.ProjectPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.ReadOnly || permission.Permission != utils.Config {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong project permission.",
		})
	}
	_, err := service.DeleteProjectPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove permission success.",
		})
	}
}

func (pc *PermissionController) AddEnvironmentToUser(ctx iris.Context) {
	var permission models.EnvironmentPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddEnvironmentPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveEnvironmentFromUser(ctx iris.Context) {
	var permission models.EnvironmentPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteEnvironmentPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}

func (pc *PermissionController) AddEnvironmentToGroup(ctx iris.Context) {
	var permission models.EnvironmentPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddEnvironmentPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveEnvironmentFromGroup(ctx iris.Context) {
	var permission models.EnvironmentPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteEnvironmentPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}

func (pc *PermissionController) AddConfigurationToUser(ctx iris.Context) {
	var permission models.ConfigurationPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddConfigurationPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveConfigurationFromUser(ctx iris.Context) {
	var permission models.ConfigurationPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteConfigurationPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}

func (pc *PermissionController) AddConfigurationToGroup(ctx iris.Context) {
	var permission models.ConfigurationPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddConfigurationPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveConfigurationFromGroup(ctx iris.Context) {
	var permission models.ConfigurationPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteConfigurationPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}

func (pc *PermissionController) AddSecretToUser(ctx iris.Context) {
	var permission models.SecretPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddSecretPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveSecretFromUser(ctx iris.Context) {
	var permission models.SecretPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteSecretPermissionForUser(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}

func (pc *PermissionController) AddSecretToGroup(ctx iris.Context) {
	var permission models.SecretPermission
	ctx.ReadJSON(&permission)
	permission.Permission = utils.Execute
	_, err := service.AddSecretPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "add environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "add environment permission success.",
		})
	}
}

func (pc *PermissionController) RemoveSecretFromGroup(ctx iris.Context) {
	var permission models.SecretPermission
	ctx.ReadJSON(&permission)
	if permission.Permission != utils.Execute {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "wrong environment permission.",
		})
	}
	_, err := service.DeleteSecretPermissionForGroup(&permission, pc.Enforcer)
	if err != nil {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: "remove environment permission failed.",
		})
	} else {
		ctx.JSON(&PermissionResp{
			StatusCode:  iris.StatusOK,
			Description: "remove environment permission success.",
		})
	}
}
