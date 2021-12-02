package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/starship-cloud/starship-iac/server/services"
)


type UserResp struct {
	StatusCode  uint
	Description string
	models.UserEntity
}

//type DeleteUserResp struct {
//	StatusCode  uint
//	Description string
//	models.UserEntity
//}
//
//type GetUserResp struct {
//	StatusCode  uint
//	Description string
//	models.UserEntity
//}

type UsersController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

type deleteUsersResp struct {
	StatusCode bool `json:"status_code"`
}

func (uc *UsersController) Get(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)

	result, err := users_service.GetUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			UserEntity:  models.UserEntity{},
		})
	} else {
		if result != nil {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				UserEntity:  *result,
			})
		}else{
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusNotFound,
				Description: "Not found",
				UserEntity:  models.UserEntity{UserId: userReq.UserId},
			})
		}
	}
}

func (uc *UsersController) Create(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	result, err := users_service.CreateUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			UserEntity:  models.UserEntity{},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "created",
			UserEntity:  *result,
		})
	}
}

func (uc *UsersController) Delete(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := users_service.DeleteUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			UserEntity:  models.UserEntity{UserId: userReq.UserId},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "deleted",
			UserEntity:  models.UserEntity{UserId: userReq.UserId},
		})
	}
}

func (uc *UsersController) Update(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := users_service.UpdateUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			UserEntity:  models.UserEntity{UserId: userReq.UserId},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "deleted",
			UserEntity:  models.UserEntity{UserId: userReq.UserId},
		})
	}
}
