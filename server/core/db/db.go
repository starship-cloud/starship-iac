package db

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

//var collection *mongo.Collection

//refactor
type DBConfig struct {
	MongoDBConnectionUri string `mapstructure:"mongodburi"`
	MongoDBName          string `mapstructure:"mongodbname"`
	MongoDBUserName      string `mapstructure:"mongodbusername"`
	MongoDBPassword      string `mapstructure:"mongodbpassword"`
	MaxConnection        int    `mapstructure:"maxconnection"`
	RootCmdLogPath       string `mapstructure:"rootcmdlogpath"`
	RootSecret           string `mapstructure:"rootsecret"`
}

type MongoDB struct {
	DBClient                   *mongo.Client
}

func NewDB (dbConfig *DBConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbConfig.MongoDBConnectionUri)
	clientOptions.SetMaxPoolSize(uint64(dbConfig.MaxConnection))
	credential := options.Credential{
		Username: dbConfig.MongoDBUserName,
		Password: dbConfig.MongoDBPassword,
	}

	clientOptions.SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("MongoDb connect success!")
	}

	return &MongoDB{
		DBClient:   client,
	}, err
}

func (d *MongoDB)Insert(collection *mongo.Collection, data interface{}) (*mongo.InsertOneResult, error) {
	objId, err := collection.InsertOne(context.TODO(), data)

	if err != nil {
		return nil, errors.WithMessage(err, "insert one item failed.")
	}
	log.Println("action->insert,objId:", objId)
	return objId, err
}

func (d *MongoDB)Delete(collection *mongo.Collection, m bson.M) (*mongo.DeleteResult, error) {
	deleteResult, err := collection.DeleteOne(context.Background(), m)
	if err != nil {
		return nil, errors.WithMessage(err, "delete one item failed.")
	}
	return deleteResult, err
}

func (d *MongoDB)UpdateOrSave(collection *mongo.Collection, target interface{}, filter bson.M) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": target}
	updateOpts := options.Update().SetUpsert(true)
	updateResult, err := collection.UpdateOne(context.Background(), filter, update, updateOpts)
	if err != nil {
		return nil, errors.WithMessage(err, "update/save one item failed.")
	}
	return updateResult, err
}

func (d *MongoDB)Update(collection *mongo.Collection, target *interface{}, filter bson.M) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": target}
	updateResult, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return nil, errors.WithMessage(err, "update one item failed.")
	}
	return updateResult, err
}

func (d *MongoDB)GetOne(collection *mongo.Collection, m bson.M) (interface{}, error) {
	var one interface{}
	err := collection.FindOne(context.Background(), m).Decode(&one)
	if err != nil {
		return nil, errors.WithMessage(err, "get one item failed.")
	}
	return one, err
}

func (d *MongoDB)GetList(collection *mongo.Collection, m bson.M) ([]*interface{}, error) {
	var list []*interface{}
	cursor, err := collection.Find(context.Background(), m)
	if err != nil {
		return nil, errors.WithMessage(err, "get list many items failed.")
	}
	err = cursor.All(context.Background(), &list)
	if err != nil {
		return nil, errors.WithMessage(err, "get list many items failed.")
	}
	_ = cursor.Close(context.Background())

	return list, err
}