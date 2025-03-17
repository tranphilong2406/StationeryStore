package api

import (
	"StoreServer/models"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateExample(c *gin.Context) {
	var req models.Example

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	example := models.NewExample(req)

	fmt.Println("example: ", example)

	res, err := models.ExampleDB.Create(example)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	response.MyResponse.Success(c, res)
}

func GetExample(c *gin.Context) {
	response.MyResponse.Error(c, myerror.EmptyParam())
}
