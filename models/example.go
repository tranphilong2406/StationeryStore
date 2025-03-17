package models

import (
	"StoreServer/job"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var ExampleDB = job.DB{
	ColName:     "example",
	DBName:      "",
	TemplateObj: Example{},
}

type Example struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	UUID        string        `bson:"uuid" json:"uuid"`
	Name        string        `bson:"name" json:"name"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

func NewExample(ex Example) *Example {
	now := time.Now()
	return &Example{
		UUID:        uuid.NewString(),
		Name:        ex.Name,
		CreatedTime: &now,
		UpdatedTime: &now,
	}
}

func InitExampleDB() {
	ExampleDB.Init("example")
}
