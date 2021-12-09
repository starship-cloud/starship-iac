package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/starship-cloud/starship-iac/server/events/models"
)

func CreateRole(role models.Role, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddPolicy(role.RoleName, role.Id, role.Permission)
}

func AddRoleForUser(roleForUser models.RoleForUser, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddRoleForUser(roleForUser.UserId, roleForUser.RoleName)
}

func DeleteRoleForUser(roleForUser models.RoleForUser, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.DeleteRoleForUser(roleForUser.UserId, roleForUser.RoleName)
}

func GetRoleForUser(userId string, enforcer *casbin.Enforcer) ([]string, error) {
	return enforcer.GetRolesForUser(userId)
}

func AddPermission(permission models.Permission, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddPolicy(permission.UserId, permission.Id, permission.Permission)
}

func DeletePermission(permission models.Permission, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.DeletePermission(permission.UserId, permission.Id, permission.Permission)
}
