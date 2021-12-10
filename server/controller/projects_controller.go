package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	"github.com/starship-cloud/starship-iac/server/services"
	"strconv"
)

type ProjectResp struct {
	StatusCode  uint
	Description string
	Data        models.ProjectEntity
}

type ProjectsResp struct {
	StatusCode  uint
	Description string
	Data        []models.ProjectEntity
}

///////

type ProjectsController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

func (pc *ProjectsController) Get(ctx iris.Context) {
	var prjReq models.ProjectEntity
	ctx.ReadJSON(&prjReq)

	userId := ctx.Params().Get("projectid")

	result, err := service.GetUserByUserId(userId, pc.DB)
	if err != nil {
		pc.Logger.Err(err.Error())
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
				Data:        models.UserEntity{UserId: prjReq.ProjectId},
			})
		}
	}
}

func (pc *ProjectsController) Create(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	result, err := service.CreateUser(&userReq, pc.DB)
	if err != nil {
		pc.Logger.Err(err.Error())
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

func (pc *ProjectsController) Delete(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := service.DeleteUser(&userReq, pc.DB)
	if err != nil {
		pc.Logger.Err(err.Error())
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

func (pc *ProjectsController) Update(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)
	_, err := service.UpdateUser(&userReq, pc.DB)
	if err != nil {
		pc.Logger.Err(err.Error())
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

func (pc *ProjectsController) Search(ctx iris.Context) {
	userName :=  ctx.URLParam("projectname")
	page_index, _ := strconv.Atoi(ctx.URLParam("page_index"))
	page_limit, _ := strconv.Atoi(ctx.URLParam("page_limit"))

	pageinOption := &models.PaginOption{
		Limit: int64(page_limit),
		Index: int64(page_index),
	}

	result, err := service.SearchUsers(userName, pc.DB, pageinOption)
	if err != nil {
		pc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{UserName: userName},
		})
	} else {
		if result != nil {
			ctx.JSON(&ProjectsResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				Data:        nil /*result*/,
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