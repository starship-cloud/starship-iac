package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/starship-cloud/starship-iac/server/services"
)

type CreateUserResp struct {
	StatusCode uint
	Reason     string
	models.UserEntity
}

type DeleteUserReq struct {
	UserId string `json:"userid"`
}

type DeleteUserResp struct {
	UserId     string `json:"userid"`
	StatusCode uint   `json:"status_code"`
}

type UsersController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

type deleteUsersResp struct {
	StatusCode bool `json:"status_code"`
}

func (uc *UsersController) Get(ctx iris.Context) {

}

func (uc *UsersController) Create(ctx iris.Context) {
	var createUserReq models.UserEntity
	ctx.ReadJSON(&createUserReq)
	result, err := users_service.CreateUser(&createUserReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&CreateUserResp{
			StatusCode: iris.StatusInternalServerError,
			Reason: err.Error(),
			UserEntity: models.UserEntity{},
		})
	} else {
		ctx.JSON(&CreateUserResp{
			StatusCode: iris.StatusConflict,
			Reason: "N/A",
			UserEntity: *result,
		})
	}
}

func (uc *UsersController) Delete(ctx iris.Context) {

}
