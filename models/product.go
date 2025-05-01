package models

import (
	"StoreServer/job"
	"StoreServer/utils/response"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ProductDB = job.DB{
	ColName:     "product",
	DBName:      "",
	TemplateObj: Product{},
}

type Product struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
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

type ProductOrder struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	Price    int    `bson:"price" json:"price"`
	Discount int    `bson:"discount" json:"discount"`
}

func (p *Product) Validate() response.Response {
	res := response.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}

	if p.Name == "" {
		res.Code = http.StatusBadRequest
		res.Message = "Name cannot be empty"
	}
	if p.Price <= 0 {
		res.Code = http.StatusBadRequest
		res.Message = "Price must be greater than 0"
	}
	if p.Stock < 0 {
		res.Code = http.StatusBadRequest
		res.Message = "Stock cannot be negative"
	}

	return res
}

type ListProduct struct {
	Products []Product `bson:"products" json:"products"`
}

func InitProductDB() {
	ProductDB.Init("store")
}
