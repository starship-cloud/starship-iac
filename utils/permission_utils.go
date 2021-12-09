package PermissionUtils

const (
	Admin         = "*"              //role1 , * , *
	ProjectCreate = "project_create" //role2 , * , project_create
	Secret        = "secret"         //role3 , * , secret

	Config   = "config"    //tom , project1 , config
	ReadOnly = "read_only" //tom , project1 , read_only
	Execute  = "execute"   //tom , sit , execute
)
