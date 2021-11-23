package db

import (
	"context"
	"fmt"
	"github.com/starship-cloud/starship-iac/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Collection

func Init() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(utils.MongoDBConnectionUri)
	clientOptions.SetMaxPoolSize(utils.MaxConnection)
	credential := options.Credential{
		Username: utils.MongoDBUserName,
		Password: utils.MongoDBPassword,
	}
	clientOptions.SetAuth(credential)
	db, err := mongo.Connect(ctx, clientOptions)
	collection = db.Database(utils.MongoDBName).Collection(utils.MongoDBName)
	if err != nil {
		fmt.Println("MongoDb connect success!")
	}
	return err
}

func Insert(data interface{}) bool {
	objId, err := collection.InsertOne(context.TODO(), data)

	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("action->insert,objId:", objId)
	return true
}

func Delete(m bson.M) bool {
	deleteResult, err := collection.DeleteOne(context.Background(), m)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("action->delete,:", deleteResult)
	return true
}

func UpdateOrSave(target interface{}, filter bson.M) bool {
	update := bson.M{"$set": target}
	updateOpts := options.Update().SetUpsert(true)
	updateResult, err := collection.UpdateOne(context.Background(), filter, update, updateOpts)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("action->update,:", updateResult)
	return true
}

func Update(target *interface{}, filter bson.M) bool {
	update := bson.M{"$set": target}
	updateResult, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("action->update,:", updateResult)
	return true
}

func GetOne(m bson.M) interface{} {
	var one interface{}
	err := collection.FindOne(context.Background(), m).Decode(&one)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("action->find one,: ", one)
	return one
}

func GetList(m bson.M) []*interface{} {
	var list []*interface{}
	cursor, err := collection.Find(context.Background(), m)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = cursor.All(context.Background(), &list)
	if err != nil {
		log.Println(err)
		return nil
	}
	_ = cursor.Close(context.Background())

	log.Println("action->find list,: ", list)
	return list
}
