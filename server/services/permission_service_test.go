package service

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"testing"
)

func Test_CreateRole(t *testing.T) {
	uri := utils.MongoDBConnectionUri
	if !strings.HasPrefix(uri, "mongodb+srv://") && !strings.HasPrefix(uri, "mongodb://") {
		uri = fmt.Sprint("mongodb://" + uri)
	}

	dbConfig := db.DBConfig{
		MongoDBConnectionUri: utils.MongoDBConnectionUri,
		MongoDBName:          utils.MongoAuthDBName,
		MongoDBUserName:      utils.MongoDBUserName,
		MongoDBPassword:      utils.MongoDBPassword,
		MaxConnection:        utils.MaxConnection,
		RootCmdLogPath:       utils.RootCmdLogPath,
		RootSecret:           utils.RootSecret,
	}
	clientOptions := options.Client().ApplyURI(dbConfig.MongoDBConnectionUri)
	clientOptions.SetMaxPoolSize(uint64(dbConfig.MaxConnection))
	credential := options.Credential{
		Username: dbConfig.MongoDBUserName,
		Password: dbConfig.MongoDBPassword,
	}

	clientOptions.SetAuth(credential)

	a, err := mongodbadapter.NewAdapterWithClientOption(clientOptions, utils.MongoAuthDBName)
	if err != nil {
		panic(err)
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "m = g(r.sub, p.sub) && ( r.obj == p.obj || p.obj==\"*\" ) && ( r.act == p.act || p.act==\"*\" )")

	e, err := casbin.NewEnforcer(m, a)
	e.EnableAutoSave(true)
	if err != nil {
		panic(err)
	}

	admin_role := &models.Role{
		RoleName:   "admin",
		Id:         "*",
		Permission: "*",
	}
	CreateRole(admin_role, e)

	projectCreater := &models.Role{
		RoleName:   "projectCreater",
		Id:         "*",
		Permission: "project_create",
	}
	CreateRole(projectCreater, e)

	security := &models.Role{
		RoleName:   "security",
		Id:         "*",
		Permission: "secret",
	}
	CreateRole(security, e)
	//name:="bob3"
	//e.AddPolicy(name, "data1", "read")
	//e.AddPolicy(name, "data1", "write")
	////e.DeletePermissionsForUser("alice")
	////e.GetNamedPolicy(name)
	//userPolicy:=e.GetFilteredPolicy(0,name)
	//fmt.Println(userPolicy)
	//
	//dataPolicy:=e.GetFilteredPolicy(1,"data1")
	//fmt.Println(dataPolicy)

}

func Test_GetAllRoles(t *testing.T) {
	uri := utils.MongoDBConnectionUri
	if !strings.HasPrefix(uri, "mongodb+srv://") && !strings.HasPrefix(uri, "mongodb://") {
		uri = fmt.Sprint("mongodb://" + uri)
	}

	dbConfig := db.DBConfig{
		MongoDBConnectionUri: utils.MongoDBConnectionUri,
		MongoDBName:          utils.MongoAuthDBName,
		MongoDBUserName:      utils.MongoDBUserName,
		MongoDBPassword:      utils.MongoDBPassword,
		MaxConnection:        utils.MaxConnection,
		RootCmdLogPath:       utils.RootCmdLogPath,
		RootSecret:           utils.RootSecret,
	}
	clientOptions := options.Client().ApplyURI(dbConfig.MongoDBConnectionUri)
	clientOptions.SetMaxPoolSize(uint64(dbConfig.MaxConnection))
	credential := options.Credential{
		Username: dbConfig.MongoDBUserName,
		Password: dbConfig.MongoDBPassword,
	}

	clientOptions.SetAuth(credential)

	a, err := mongodbadapter.NewAdapterWithClientOption(clientOptions, utils.MongoAuthDBName)
	if err != nil {
		panic(err)
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "m = g(r.sub, p.sub) && ( r.obj == p.obj || p.obj==\"*\" ) && ( r.act == p.act || p.act==\"*\" )")

	e, err := casbin.NewEnforcer(m, a)
	e.EnableAutoSave(true)
	if err != nil {
		panic(err)
	}

	fmt.Println(GetAllRoles(e))

}
