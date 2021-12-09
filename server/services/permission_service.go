package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/starship-cloud/starship-iac/server/events/models"
)

func AddPermission(permission models.Permission, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddPolicy(permission.UserId, permission.ProjectId, permission.Permission)
}
