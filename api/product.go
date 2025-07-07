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

	fmt.Println("req: ", req)

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	res := models.ProductDB.Create(req)

	c.JSON(res.Code, res)
}

func GetProduct(c *gin.Context) {
	page := utils.ParseInt(c.Query("page"), 1)
	pageSize := utils.ParseInt(c.Query("page_size"), 10)
	category_id := c.Query("category_id")
	name := c.Query("search")
	var filter bson.M

	if category_id == "" {
		if name != "" {
			filter = bson.M{
				"deleted_time": nil,
				"name":         bson.M{"$regex": name, "$options": "i"},
			}
		} else {
			filter = bson.M{
				"deleted_time": nil,
			}
		}
	} else {
		if name != "" {
			filter = bson.M{
				"deleted_time": nil,
				"category_id":  category_id,
				"name":         bson.M{"$regex": name, "$options": "i"},
			}
		} else {
			filter = bson.M{
				"deleted_time": nil,
				"category_id":  category_id,
			}
		}
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

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	filter := bson.M{
		"_id":          objID,
		"deleted_time": nil,
	}

	res := models.ProductDB.QueryOne(filter)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	update := res.Data.(*models.Product)

	if req.Image != "" {
		update.Image = req.Image
	}

	update.Name = req.Name
	update.Description = req.Description

	update.Stock = req.Stock
	update.SellPrice = req.SellPrice
	update.BuyPrice = req.BuyPrice

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
