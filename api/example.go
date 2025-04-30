package api

import (
	"StoreServer/models"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateExample(c *gin.Context) {
	var req models.Example

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
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

	for _, v := range req {
		if ok := v.Validate(); ok.Code != http.StatusOK {
			c.JSON(ok.Code, ok)
			return
		}
	}

	lst := make([]interface{}, len(req))
	for i, v := range req {
		lst[i] = v
	}

	res := models.ExampleDB.CreateMany(lst)

	c.JSON(res.Code, res)
}

func GetExample(c *gin.Context) {
	response.MyResponse.Error(c, myerror.EmptyParam())
}

func UpdateExample(c *gin.Context) {
	var req models.Example

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

	res := models.ExampleDB.QueryOne(filter)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	update := res.Data.(*models.Example)

	update.Name = req.Name

	updating := models.ExampleDB.UpdateOne(filter, update)

	c.JSON(updating.Code, updating)
}

func DeleteExample(c *gin.Context) {
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

	res := models.ExampleDB.DeleteOne(filter)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	c.JSON(res.Code, res)
}
