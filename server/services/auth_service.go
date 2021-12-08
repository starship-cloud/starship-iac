package service

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/starship-cloud/starship-iac/utils"
	"time"
)

func CreateToken(userId string) (string, error) {
	now := time.Now()
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iat":    now.Unix(),
		"exp":    now.Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(utils.RootSecret))
}
