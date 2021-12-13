package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/starship-cloud/starship-iac/server/events/models"
)

func CreateRole(role *models.Role, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddPolicy(role.RoleName, role.Id, role.Permission)
}

func AddRoleForUser(roleForUser *models.RoleForUser, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddRoleForUser(roleForUser.UserId, roleForUser.RoleName)
}

func DeleteRoleForUser(roleForUser *models.RoleForUser, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.DeleteRoleForUser(roleForUser.UserId, roleForUser.RoleName)
}

func GetRoleForUser(userId string, enforcer *casbin.Enforcer) ([]string, error) {
	return enforcer.GetRolesForUser(userId)
}

func AddProjectPermissionForUser(permission *models.ProjectPermission, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.AddPolicy(permission.UserId, permission.ProjectId, permission.Permission)
}

func DeleteProjectPermissionForUser(permission *models.ProjectPermission, enforcer *casbin.Enforcer) (bool, error) {
	return enforcer.RemovePolicy(permission.UserId, permission.ProjectId, permission.Permission)
}

func GetAllProjectPermissionsForUser(userId string, enforcer *casbin.Enforcer) [][]string {
	return enforcer.GetFilteredPolicy(0, userId)
}

func GetUsersByProjectId(projectId string, enforcer *casbin.Enforcer) [][]string {
	return enforcer.GetFilteredPolicy(1, projectId)
}
