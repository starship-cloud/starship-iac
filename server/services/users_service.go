package users_service

import (
	"fmt"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DB_NAME       = "starship"
	DB_COLLECTION = "users"
)

func GetUser(userId string, db *db.MongoDB) (*models.UserEntity, error){
	collection := db.DBClient.Database(DB_NAME).Collection(DB_COLLECTION)

	filter := bson.M{"userid": userId}
	result, err := db.GetOne(collection, filter)

	if err != nil {
		return nil, fmt.Errorf("get user with user id %s failed due to DB operation", userId)
	}else if result != nil{
		return result.(*models.UserEntity), nil
	}else{
		//not found
		return nil, nil
	}
}


func CreateUser(user *models.UserEntity, db *db.MongoDB) (*models.UserEntity ,error){
	collection := db.DBClient.Database(DB_NAME).Collection(DB_COLLECTION)

	filter := bson.M{"username": user.Username}
	result, err := db.GetOne(collection, filter)

	if err != nil {
		return nil, err
	}else if result != nil{
		return nil, fmt.Errorf("the user %s has been exist.", user.Username)
	}else{
		userId := utils.GenUserId()
		newUser := &models.UserEntity{
			Userid:     userId,
			Username:   user.Username,
			Email:      user.Email,
			Password:   user.Password,
		}

		_, err := db.UpdateOrSave(collection, newUser, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("update/save user %s failed due to DB operation", user.Username)
		}else{
			return newUser, nil
		}
	}
}

func DeleteUser(user *models.UserEntity, db *db.MongoDB) (*models.UserEntity ,error){
	collection := db.DBClient.Database(DB_NAME).Collection(DB_COLLECTION)

	filter := bson.M{"userid": user.Userid}
	_, err := db.Delete(collection, filter)

	if err != nil {
		return nil, fmt.Errorf("the user with user id: %s has been deleted.", user.Userid)
	}else{
		return nil, nil
	}
}

func UpdateUser(user *models.UserEntity, db *db.MongoDB) (*models.UserEntity ,error){
	collection := db.DBClient.Database(DB_NAME).Collection(DB_COLLECTION)

	filter := bson.M{"userid": user.Userid}
	result, err := db.GetOne(collection, filter)

	if err != nil {
		return nil, err
	}else if result != nil{
		newUser := &models.UserEntity{
			Userid:     result.(*models.UserEntity).Userid,
			Username:   user.Username,
			Email:      user.Email,
			Password:   user.Password,
		}

		_, err := db.UpdateOrSave(collection, newUser, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("update/save user %s failed due to DB operation", user.Username)
		}else{
			return newUser, nil
		}
	}else{
		return nil, fmt.Errorf("the user wit user id %s not exist.", user.Userid)
	}
}