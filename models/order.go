package models

import (
	"StoreServer/job"
	"StoreServer/utils/response"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var OrderDB = job.DB{
	ColName:     "order",
	DBName:      "",
	TemplateObj: Order{},
}

type Order struct {
	ID          bson.ObjectID  `json:"id" bson:"_id,omitempty"`
	Description string         `json:"description" bson:"description"`
	Products    []ProductOrder `json:"products" bson:"products" binding:"required"`
	TotalPrice  float64        `json:"total_price" bson:"total_price" binding:"required"`
	Discount    float64        `json:"discount" bson:"discount"`
	Status      bool           `json:"status" bson:"status"` // 0 is pay 1 is not pay
	CreatedTime *time.Time     `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time     `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time     `json:"deleted_time" bson:"deleted_time,omitempty"`
}

type UpdateOrder struct {
	ID     bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Status bool          `json:"status" bson:"status"`
}

func (o *Order) Validate() response.Response {
	res := response.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}

	if len(o.Products) < 1 {
		res.Code = http.StatusBadRequest
		res.Message = "At least 1 product in order"
	}
	if o.TotalPrice <= 0 {
		res.Code = http.StatusBadRequest
		res.Message = "Price must be greater than 0"
	}

	return res
}

func InitOrderDB() {
	OrderDB.Init("store")
}
