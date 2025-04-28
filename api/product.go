package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
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
	filter := models.Product{}
	filter.DeletedTime = nil

	offset := (page - 1) * pageSize
	res := models.ProductDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.(models.ListProduct)

	c.JSON(res.Code, res)
}

func UpdateProduct() {}

func DeleteProduct() {}
