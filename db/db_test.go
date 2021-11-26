package db

import (
	"fmt"
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
	err := Init()
	if err != nil {
		return
	}

	data := Student{
		Name:   "tom",
		Age:    18,
		Sid:    "001",
		Status: 1,
	}
	result := Insert(&data)
	Assert(t, result != true, "result should be true")
}

func Test_GetOne(t *testing.T) {
	err := Init()
	if err != nil {
		return
	}

	filter := bson.M{"name": "tom"}
	result := GetOne(filter)
	fmt.Println(result)

	Assert(t, result != nil, "result should be nil")
}

func Test_UpdateOrSave(t *testing.T) {
	err := Init()
	if err != nil {
		return
	}

	data := Student{
		Name:   "jerry",
		Age:    20,
		Sid:    "002",
		Status: 2,
	}
	result := UpdateOrSave(&data, bson.M{})
	fmt.Println(result)
	Assert(t, result != true, "result should be true")
}

func Test_GetList(t *testing.T) {
	err := Init()
	if err != nil {
		return
	}
	result := GetList(bson.M{})
	fmt.Println(result)
}

func Test_Delete(t *testing.T) {
	err := Init()
	if err != nil {
		return
	}
	filter := bson.M{"name": "tom"}

	result := Delete(filter)
	fmt.Println(result)
	Assert(t, result, "result should be true")
}
