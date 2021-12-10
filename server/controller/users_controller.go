package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/starship-cloud/starship-iac/server/services"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UserResp struct {
	StatusCode  uint
	Description string
	Data        models.UserEntity
}

type UsersResp struct {
	StatusCode  uint
	Description string
	Data        []models.UserEntity
}

///////

type UsersController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

func (uc *UsersController) Login(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)

	user, err := service.GetUserByNmae(userReq.UserName, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{},
		})
	} else {
		if user != nil {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
			if err != nil {
				ctx.JSON( &AuthResp{
					StatusCode: iris.StatusUnauthorized,
					Description: "password is not correct",
				})
			}

			token, _ := service.CreateToken(user.UserId)

			ctx.JSON(&AuthResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				Data: models.AuthEntity{UserId: user.UserId, AuthToken: token},
			})
		} else {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusNotFound,
				Description: "user not found",
				Data:        models.UserEntity{UserId: userReq.UserId},
			})
		}
	}
}

func (uc *UsersController) Get(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)

	userId := ctx.Params().Get("userid")

	result, err := service.GetUserByUserId(userId, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{},
		})
	} else {
		if result != nil {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				Data:        *result,
			})
		} else {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusNotFound,
				Description: "Not found",
				Data:        models.UserEntity{UserId: userReq.UserId},
			})
		}
	}
}

func (uc *UsersController) Create(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	result, err := service.CreateUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{UserName: userReq.UserName},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "created",
			Data:        *result,
		})
	}
}

func (uc *UsersController) Delete(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := service.DeleteUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{UserId: userReq.UserId},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "deleted",
			Data:        models.UserEntity{UserId: userReq.UserId},
		})
	}
}

func (uc *UsersController) Update(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := service.UpdateUser(&userReq, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{UserId: userReq.UserId},
		})
	} else {
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusOK,
			Description: "updated",
			Data:        models.UserEntity{UserId: userReq.UserId},
		})
	}
}

func (uc *UsersController) Search(ctx iris.Context) {
	userName :=  ctx.URLParam("username")
	page_index, _ := strconv.Atoi(ctx.URLParam("page_index"))
	page_limit, _ := strconv.Atoi(ctx.URLParam("page_limit"))

	pageinOption := &models.PaginOption{
		Limit: int64(page_limit),
		Index: int64(page_index),
	}

	result, err := service.SearchUsers(userName, uc.DB, pageinOption)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{UserName: userName},
		})
	} else {
		if result != nil {
			ctx.JSON(&UsersResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				Data:        result,
			})
		} else {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusNotFound,
				Description: "Not found",
				Data:        models.UserEntity{UserName: userName},
			})
		}
	}
}


