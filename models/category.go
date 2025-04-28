package models

import (
	"StoreServer/job"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var CategoryDB = job.DB{
	ColName:     "category",
	DBName:      "",
	TemplateObj: Category{},
}

type Category struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

func InitCategoryDB() {
	CategoryDB.Init("category")
}
