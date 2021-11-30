package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/logging"
)

const(
	DB_COLLECTION    =  "users"
)

type createUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type deleteUserReq struct {
	UserId string `json:"userid"`
}

type UsersController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

type createUsersResp struct {
	StatusCode  uint `json:"status_code"`
}

type deleteUsersResp struct {
	StatusCode  bool `json:"status_code"`
}

func (uc *UsersController) Get(ctx iris.Context) {
}

func (uc *UsersController) Create(ctx iris.Context) {
	data, err := json.MarshalIndent(&createUsersResp{
		StatusCode:  iris.StatusOK,
	}, "", "  ")

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		uc.Logger.Info(fmt.Sprintf("Error creating user json response: %s", err))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(data)
}

func (uc *UsersController) Delete(ctx iris.Context) {

}
