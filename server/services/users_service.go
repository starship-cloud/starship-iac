package users_service

import (
	"fmt"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DB_NAME       = "starship-db"
	DB_COLLECTION = "users"
)

func CreateUser(user *models.UserEntity, db *db.MongoDB) (*models.UserEntity ,error){
	collection := db.DBClient.Database(DB_NAME).Collection(DB_COLLECTION)

	filter := bson.M{"username": user.Username}
	result := db.GetOne(collection, filter)

	if result != nil {
		return nil, fmt.Errorf("the user %s has been exist.", user.Username)
	}else{
		userId := utils.GenUserId()
		newUser := &models.UserEntity{
			UserId:     userId,
			Username:   user.Username,
			Email:      user.Email,
			Password:   user.Password,
		}

		result := db.UpdateOrSave(collection, newUser, bson.M{})
		if result {
			return newUser, nil
		}else{
			return nil, fmt.Errorf("update/save user %s failed due to DB operation", user.Username)
		}
	}
}