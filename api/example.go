package api

import (
	"StoreServer/models"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateExample(c *gin.Context) {
	var req models.Example

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	res := models.ExampleDB.Create(req)
	if res.Code != http.StatusOK {
		//response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		c.JSON(res.Code, res)
		return
	}

	c.JSON(res.Code, res)
}

func CreateListExample(c *gin.Context) {
	var req []models.Example

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	lst := make([]interface{}, len(req))
	for i, v := range req {
		lst[i] = v
	}

	res, err := models.ExampleDB.CreateMany(lst)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	response.MyResponse.Success(c, res)
}

func GetExample(c *gin.Context) {
	response.MyResponse.Error(c, myerror.EmptyParam())
}
