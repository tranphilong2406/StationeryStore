package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
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

	res := models.ProductDB.Create(req)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}
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

func UpdateProduct() {}

func DeleteProduct() {}
