package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
)

func CreateReceivedOrder(c *gin.Context) {
	var req models.ReceivedOrder
	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	for idx := range len(req.Products) {
		filter := bson.M{
			"_id": req.Products[idx].ID,
		}

		prod := models.ProductDB.QueryOne(filter)
		if prod.Code != http.StatusOK {
			c.JSON(prod.Code, prod)
			return
		}

		update := prod.Data.(*models.Product)

		update.Stock += req.Products[idx].Quantity
		update.Price = req.Products[idx].Price

		updated := models.ProductDB.UpdateOne(filter, update)
		if updated.Code != http.StatusOK {
			c.JSON(updated.Code, updated)
			return
		}
	}

	res := models.ReceivedOrderDB.Create(req)
	c.JSON(res.Code, res)
}

func GetReceivedOrder(c *gin.Context) {
	page := utils.ParseInt(c.Query("page"), 1)
	pageSize := utils.ParseInt(c.Query("page_size"), 10)

	filter := bson.M{
		"deleted_time": nil,
	}

	offset := (page - 1) * pageSize
	res := models.ReceivedOrderDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.([]models.ReceivedOrder)
	res.Page = page
	res.PageSize = pageSize

	c.JSON(res.Code, res)
}

func GetReceivedOrderByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.MyResponse.Error(c, myerror.EmptyParam())
		return
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	result := models.ReceivedOrderDB.QueryOne(filter)

	c.JSON(result.Code, result)
}

func DeleteReceivedOrder(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.MyResponse.Error(c, myerror.EmptyParam())
		return
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	res := models.ReceivedOrderDB.DeleteOne(filter)

	c.JSON(res.Code, res)
}
