package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"github.com/starship-cloud/starship-iac/server/events/models"
	"github.com/starship-cloud/starship-iac/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

func GetProjectByNmae(userName string, db *db.MongoDB) (*models.UserEntity, error) {
	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)

	filter := bson.M{"username": userName}
	userEntity := &models.UserEntity{}
	err := db.GetOne(collection, filter, &userEntity)

	if err != nil {
		return nil, fmt.Errorf("get user %s failed due to DB operation", userName)
	} else if userEntity.UserId != "" {
		return userEntity, nil
	} else {
		//not found
		return nil, nil
	}
}

func GetProjectByProjectId(projectId string, db *db.MongoDB) (*models.ProjectEntity, error) {
	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)

	filter := bson.M{"project_id": projectId}

	pojectEntity := &models.ProjectEntity{}
	err := db.GetOne(collection, filter, &pojectEntity)

	if err != nil {
		return nil, fmt.Errorf("get project with project id %s failed due to DB operation", projectId)
	} else if pojectEntity.ProjectId != "" {
		return pojectEntity, nil
	} else {
		//not found
		return nil, nil
	}
}

func CreateProject(project *models.ProjectEntity, db *db.MongoDB) (*models.ProjectEntity, error) {
	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)
	projectEntity := &models.ProjectEntity{}

	filter := bson.M{"project_name": project.ProjectName}
	db.GetOne(collection, filter, projectEntity)

	if projectEntity.ProjectId != "" {
		return nil, fmt.Errorf("the project %s with projectId %s has been exist.", project.ProjectName, project.ProjectId)
	} else {
		projectId := utils.GenProjectId()
		t := time.Now().Unix()

		newProject := &models.ProjectEntity{
			ProjectId:    projectId,
			ProjectName:  project.ProjectName,
			Discription:  project.Discription,
			CreateAt:  t,
			UpdateAt:  t,
		}

		_, err := db.Insert(collection, newProject)
		if err != nil{
			return nil, fmt.Errorf("save project %s failed due to DB operation", project.ProjectName)
		} else{
			return newProject, nil
		}

	}
}

func UpdateProject(project *models.ProjectEntity, db *db.MongoDB) (*models.ProjectEntity, error) {
	if len(strings.TrimSpace(project.ProjectId)) == 0 ||
		len(strings.TrimSpace(project.ProjectName)) == 0  {
		return nil, errors.New("userid/username/email are required.")
	}

	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)
	projectEntity := &models.ProjectEntity{}
	filter := bson.M{"project_id": project.ProjectId}

	db.GetOne(collection, filter, projectEntity)

	if projectEntity.ProjectId != "" {
		//found
		newProject := &models.ProjectEntity{
			ProjectId:   projectEntity.ProjectId,
			ProjectName: project.ProjectName,
			Discription: project.Discription,
			CreateAt: time.Now().Unix(),
		}

		_, err := db.UpdateOrSave(collection, newProject, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("update project %s failed due to DB operation", project.ProjectName)
		} else {
			return newProject, nil
		}
	} else {
		return nil, fmt.Errorf("the user %s with user id %s not exist.", projectEntity.ProjectName, projectEntity.ProjectId)
	}
}

func DeleteProject(project *models.ProjectEntity, db *db.MongoDB) (*models.ProjectEntity, error) {
	if len(strings.TrimSpace(project.ProjectId)) == 0 {
		return nil, errors.New("projectId is required.")
	}

	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)

	projectEntity := &models.ProjectEntity{}
	filter := bson.M{"project_id": project.ProjectId}
	err := db.GetOne(collection, filter, projectEntity)

	if err != nil {
		return nil, errors.Wrap(err, "delete failed")
	} else if projectEntity.ProjectId != "" {
		//found, delete
		_, err := db.Delete(collection, filter)
		return nil, err
	} else {
		return nil, fmt.Errorf("the project with project id: %s has been deleted.", project.ProjectId)
	}

}

func SearchProjects(projectName string, db *db.MongoDB, pageinOpt *models.PaginOption) ([]models.ProjectEntity, error) {
	collection := db.DBClient.Database(models.DB_NAME).Collection(models.DB_COLLECTION_USERS)
	var projects []models.ProjectEntity
	filter := bson.M{
		"username": bson.M{
			"$regex":  projectName,
			"$options": "i",
		},
	}

	db.GetList(collection, filter, &projects, *pageinOpt)

	if len(projects) == 0 {
		return nil, fmt.Errorf("get project %s failed due to DB operation", projectName)
	} else {

		return projects, nil
	}
}
