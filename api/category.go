package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
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
	filter := models.Category{}
	filter.DeletedTime = nil

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

	//res := models.CategoryDB.Update(req)
	//if res.Code != http.StatusOK {
	//	c.JSON(res.Code, res)
	//	return
	//}
	//c.JSON(res.Code, res)
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
