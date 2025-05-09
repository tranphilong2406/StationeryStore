package models

import (
	"StoreServer/job"
	"StoreServer/utils/response"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"time"
)

var ReceivedOrderDB = job.DB{
	ColName:     "received_order",
	DBName:      "",
	TemplateObj: ReceivedOrder{},
}

type ReceivedOrder struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Products    []ProductOrder `bson:"products" json:"products"`
	TotalPrice  float64        `bson:"total_price" json:"total_price"`
	CreatedTime *time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime *time.Time     `bson:"updated_time" json:"updated_time"`
	DeletedTime *time.Time     `bson:"deleted_time" json:"deleted_time,omitempty"`
}

func (r ReceivedOrder) Validate() response.Response {
	res := response.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}

	if len(r.Products) < 1 {
		res.Code = http.StatusBadRequest
		res.Message = "At least 1 product in order"
	}
	if r.TotalPrice <= 0 {
		res.Code = http.StatusBadRequest
		res.Message = "Price must be greater than 0"
	}

	return res
}

func InitReceivedOrderDB() {
	ReceivedOrderDB.Init("store")
}
