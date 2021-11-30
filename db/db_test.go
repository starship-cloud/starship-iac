package db

import (
	"fmt"
	"github.com/starship-cloud/starship-iac/server/core/db"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func errLog(tb testing.TB, fmt string, args ...interface{}) {
	tb.Helper()
	tb.Logf("\033[31m"+fmt+"\033[39m", args...)
}

func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		errLog(tb, msg, v...)
		tb.FailNow()
	}
}

type Student struct {
	Name   string
	Age    int
	Sid    string
	Status int
}

func Test_Insert(t *testing.T) {
	err := db.Init()
	if err != nil {
		return
	}

	data := Student{
		Name:   "tom",
		Age:    18,
		Sid:    "001",
		Status: 1,
	}
	result := db.Insert(&data)
	Assert(t, result != true, "result should be true")
}

func Test_GetOne(t *testing.T) {
	err := db.Init()
	if err != nil {
		return
	}

	filter := bson.M{"name": "tom"}
	result := db.GetOne(filter)
	fmt.Println(result)

	Assert(t, result != nil, "result should be nil")
}

func Test_UpdateOrSave(t *testing.T) {
	err := db.Init()
	if err != nil {
		return
	}

	data := Student{
		Name:   "jerry",
		Age:    20,
		Sid:    "002",
		Status: 2,
	}
	result := db.UpdateOrSave(&data, bson.M{})
	fmt.Println(result)
	Assert(t, result != true, "result should be true")
}

func Test_GetList(t *testing.T) {
	err := db.Init()
	if err != nil {
		return
	}
	result := db.GetList(bson.M{})
	fmt.Println(result)
}

func Test_Delete(t *testing.T) {
	err := db.Init()
	if err != nil {
		return
	}
	filter := bson.M{"name": "tom"}

	result := db.Delete(filter)
	fmt.Println(result)
	Assert(t, result, "result should be true")
}
