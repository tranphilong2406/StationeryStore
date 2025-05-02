package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateProduct(c *gin.Context) {
	var req models.Product

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	res := models.ProductDB.Create(req)

	c.JSON(res.Code, res)
}

func GetProduct(c *gin.Context) {
	page := utils.ParseInt(c.Param("page"), 1)
	pageSize := utils.ParseInt(c.Param("page_size"), 10)
	filter := bson.M{
		"deleted_time": nil,
	}

	offset := (page - 1) * pageSize
	res := models.ProductDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.([]models.Product)

	c.JSON(res.Code, res)
}

func UpdateProduct(c *gin.Context) {
	var req models.Product

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	filter := bson.M{
		"_id":          req.ID,
		"deleted_time": nil,
	}

	res := models.ProductDB.QueryOne(filter)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	update := res.Data.(*models.Product)
	fmt.Println("update: ", update)

	update.Name = req.Name
	update.Description = req.Description
	update.Image = req.Image
	update.Stock = req.Stock
	update.Price = req.Price

	updating := models.ProductDB.UpdateOne(filter, update)

	c.JSON(updating.Code, updating)
}

func DeleteProduct(c *gin.Context) {
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

	res := models.ProductDB.DeleteOne(filter)

	c.JSON(res.Code, res)
}
