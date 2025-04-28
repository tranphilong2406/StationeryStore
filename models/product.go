package models

import (
	"StoreServer/job"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var ProductDB = job.DB{
	ColName:     "product",
	DBName:      "",
	TemplateObj: Product{},
}

type Product struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	CategoryID  string        `bson:"category_id" json:"category_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Image       string        `bson:"image" json:"image"`
	Stock       int           `bson:"stock" json:"stock"`
	Price       int           `bson:"price" json:"price"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

type ListProduct struct {
	Products []Product `bson:"products" json:"products"`
}

func InitProductDB() {
	ProductDB.Init("product")
}
