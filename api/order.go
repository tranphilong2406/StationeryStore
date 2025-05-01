package api

import (
	"StoreServer/models"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateOrder(c *gin.Context) {
	var req models.Order
	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	for idx := range len(req.Products) {
		temp, err := bson.ObjectIDFromHex(req.Products[idx].ID)
		if err != nil {
			response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
			return
		}

		filter := bson.M{
			"_id": temp,
		}

		prod := models.ProductDB.QueryOne(filter)
		if prod.Code != http.StatusOK {
			c.JSON(prod.Code, prod)
			return
		}

		update := prod.Data.(*models.Product)

		update.Stock -= req.Products[idx].Quantity

		updated := models.ProductDB.UpdateOne(filter, update)
		if updated.Code != http.StatusOK {
			c.JSON(updated.Code, updated)
			return
		}
	}

	res := models.OrderDB.Create(req)
	c.JSON(res.Code, res)
}

func GetOrder(c *gin.Context) {
	// Implementation for getting an order
}

func UpdateOrder(c *gin.Context) {
	// Implementation for updating an order
}

func DeleteOrder(c *gin.Context) {
	// Implementation for deleting an order
}
