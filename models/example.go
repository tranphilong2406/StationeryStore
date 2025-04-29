package models

import (
	"StoreServer/job"
	"StoreServer/utils/response"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ExampleDB = job.DB{
	ColName:     "example",
	DBName:      "",
	TemplateObj: Example{},
}

type Example struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

func (e Example) Validate() response.Response {
	res := response.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}

	if e.Name == "" {
		res.Code = http.StatusBadRequest
		res.Message = "Name cannot be empty"
	}

	return res
}

type ListExample struct {
	Examples []Example `bson:"examples" json:"examples"`
}

func InitExampleDB() {
	ExampleDB.Init("store")
}
