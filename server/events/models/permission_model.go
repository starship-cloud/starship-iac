package models

type ProjectPermissionRequestBody struct {
	Id          string
	IsGroup     bool
	ProjectId   string
	IsAddAction bool
}

type ProjectPermission struct {
	UserId     string
	GroupId    string
	ProjectId  string
	Permission string
}

type Role struct {
	RoleName   string
	Id         string
	Permission string
}

type RoleForUser struct {
	RoleName string
	UserId   string
}

type EnvironmentPermission struct {
	UserId        string
	GroupId       string
	EnvironmentId string
	Permission    string
}
