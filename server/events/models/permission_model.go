package models

type ProjectPermission struct {
	Id         string
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
	Id            string
	EnvironmentId string
	Permission    string
}

type ConfigurationPermission struct {
	Id              string
	ConfigurationId string
	Permission      string
}

type SecretPermission struct {
	Id         string
	SecretId   string
	Permission string
}
