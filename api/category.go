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

func CreateCategory(c *gin.Context) {
	var req models.Category

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	res := models.CategoryDB.Create(req)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}
	c.JSON(res.Code, res)
}

func GetCategory(c *gin.Context) {
	page := utils.ParseInt(c.Param("page"), 1)
	pageSize := utils.ParseInt(c.Param("page_size"), 10)
	filter := bson.M{
		"deleted_time": nil,
	}

	offset := (page - 1) * pageSize
	res := models.CategoryDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	c.JSON(res.Code, res)
}

func UpdateCategory(c *gin.Context) {
	var req models.Category

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

	res := models.CategoryDB.QueryOne(filter)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	update := res.Data.(*models.Category)

	update.Name = req.Name
	update.Description = req.Description

	updating := models.CategoryDB.Update(filter, update)

	c.JSON(updating.Code, updating)
}

func DeleteCategory(c *gin.Context) {
	var req models.Category

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	//res := models.CategoryDB.Delete(req)
	//if res.Code != http.StatusOK {
	//	c.JSON(res.Code, res)
	//	return
	//}
	//c.JSON(res.Code, res)
}
