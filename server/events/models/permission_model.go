package models

type Permission struct {
	UserId     string
	ProjectId  string
	Permission string
}

type Role struct {
	RoleName   string
	ProjectId  string
	Permission string
}
