package models

type ProjectPermission struct {
	UserId     string
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
