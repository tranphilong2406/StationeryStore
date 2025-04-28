package models

import (
	"StoreServer/job"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ExampleDB = job.DB{
	ColName:     "example",
	DBName:      "",
	TemplateObj: Example{},
}

type Example struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

type ListExample struct {
	Examples []Example `bson:"examples" json:"examples"`
}

func InitExampleDB() {
	ExampleDB.Init("store")
}
