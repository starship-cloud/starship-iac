package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/server/logging"
	service "github.com/starship-cloud/starship-iac/server/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthResp struct {
	StatusCode  uint
	Description string
	Data        models.AuthEntity
}

type AuthController struct {
	Logger  logging.SimpleLogging
	Drainer *events.Drainer
	DB      *db.MongoDB
}

func (uc *AuthController) Login(ctx iris.Context) {
	var userReq models.UserEntity
	ctx.ReadJSON(&userReq)

	user, err := service.GetUserByNmae(userReq.Username, uc.DB)
	if err != nil {
		uc.Logger.Err(err.Error())
		ctx.JSON(&UserResp{
			StatusCode:  iris.StatusInternalServerError,
			Description: err.Error(),
			Data:        models.UserEntity{},
		})
	} else {
		if user != nil {
			//found
			//compare hash by Hashed(password) == user.Password
			// if not return AuthResp with 401
			// else create auth token and return 200 with an AuthResp
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
			if err != nil {
				ctx.JSON( &AuthResp{
					StatusCode: iris.StatusUnauthorized,
					Description: "password is not correct",
				})
			}

			token, _ := service.CreateToken(user.Userid)

			ctx.JSON(&AuthResp{
				StatusCode:  iris.StatusOK,
				Description: "found",
				Data: models.AuthEntity{UserId: user.Userid, AuthToken: token},
			})
		} else {
			ctx.JSON(&UserResp{
				StatusCode:  iris.StatusNotFound,
				Description: "user not found",
				Data:        models.UserEntity{Userid: userReq.Userid},
			})
		}
	}
}