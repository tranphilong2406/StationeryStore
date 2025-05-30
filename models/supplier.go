package models

import (
	"StoreServer/job"
	"StoreServer/utils/response"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var SupplierDB = job.DB{
	ColName:     "supplier",
	DBName:      "",
	TemplateObj: Supplier{},
}

type Supplier struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Phone       string        `bson:"phone" json:"phone"`
	Email       string        `bson:"email" json:"email"`
	Address     string        `bson:"address" json:"address"`
	CreatedTime *time.Time    `json:"created_time" bson:"created_time"`
	UpdatedTime *time.Time    `json:"updated_time" bson:"updated_time"`
	DeletedTime *time.Time    `json:"deleted_time" bson:"deleted_time,omitempty"`
}

func (s Supplier) Validate() response.Response {
	res := response.Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}

	if s.Name == "" {
		res.Code = http.StatusBadRequest
		res.Message = "Name cannot be empty"
	}

	if s.Phone == "" {
		res.Code = http.StatusBadRequest
		res.Message = "Phone cannot be empty"
	}

	if s.Email == "" {
		res.Code = http.StatusBadRequest
		res.Message = "Email cannot be empty"
	}

	return res
}

func InitSupplierDB() {
	SupplierDB.Init("store")
}
