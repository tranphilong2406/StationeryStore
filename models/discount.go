package models

import (
	"StoreServer/job"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DiscountDB = job.DB{
	ColName:     "discount",
	DBName:      "",
	TemplateObj: Discount{},
}

type Discount struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Code                string        `bson:"code" json:"code"`
	Description         string        `bson:"description" json:"description"`
	Type                string        `bson:"type" json:"type"`
	Value               int           `bson:"value" json:"value"`
	StartDateTime       string        `bson:"start_date_time" json:"start_date_time"`
	EndDateTime         string        `bson:"end_date_time" json:"end_date_time"`
	Status              string        `bson:"status" json:"status"`
	MinPurchase         int           `bson:"min_purchase" json:"min_purchase"`
	ApplicableProductID string        `bson:"application_product_id" json:"application_product_id"`
	MaxUsage            int           `bson:"max_usage" json:"max_usage"`
	UsedCount           int           `bson:"used_count" json:"used_count"`
	CreatedTime         string        `bson:"created_time" json:"created_time"`
	UpdatedTime         string        `bson:"updated_time" json:"updated_time"`
	DeletedTime         string        `bson:"deleted_time" json:"deleted_time,omitempty"`
}
