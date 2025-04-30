package response

import (
	myerr "StoreServer/utils/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Code     int         `json:"code"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
}

var MyResponse Response

func Respond(c *gin.Context, response *Response) {
	c.JSON(response.Code, response)
}

func (Response) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    data,
	})
}

func (Response) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "ok",
		"data":    data,
	})
}

func (Response) Error(c *gin.Context, err myerr.CustomError) {
	c.JSON(err.HTTPCode, gin.H{
		"code":    err.HTTPCode,
		"message": err.Error(),
	})
}
