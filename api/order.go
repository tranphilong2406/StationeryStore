package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	gettime "StoreServer/utils/get_time"
	"StoreServer/utils/response"
	"fmt"
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
		filter := bson.M{
			"_id": req.Products[idx].ID,
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
	page := utils.ParseInt(c.Query("page"), 1)
	pageSize := utils.ParseInt(c.Query("page_size"), 10)
	status := c.Query("status")

	timeStart, timeEnd := gettime.RangeFromKeyword(c.Query("time_range"))

	filter := bson.M{
		"deleted_time": nil,
		"created_time": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
	}

	if status != "all" {
		filter["status"] = utils.ParseBool(status)
	}

	fmt.Println("filter:", filter)

	offset := (page - 1) * pageSize
	res := models.OrderDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.([]models.Order)
	res.Page = page
	res.PageSize = pageSize

	c.JSON(res.Code, res)
}

func UpdateStatusOrder(c *gin.Context) {
	var req models.UpdateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	filter := bson.M{"_id": req.ID}
	update := bson.M{"status": req.Status}

	res := models.OrderDB.UpdateOne(filter, update)
	c.JSON(res.Code, res)
}

func GetOrderByID(c *gin.Context) {
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

	result := models.OrderDB.QueryOne(filter)

	c.JSON(result.Code, result)
}

func DeleteOrder(c *gin.Context) {
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

	res := models.OrderDB.DeleteOne(filter)

	c.JSON(res.Code, res)
}
