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

func CreateSupplier(c *gin.Context) {
	var req models.Supplier

	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	res := models.SupplierDB.Create(req)

	c.JSON(res.Code, res)
}

func GetSupplier(c *gin.Context) {
	page := utils.ParseInt(c.Query("page"), 1)
	pageSize := utils.ParseInt(c.Query("page_size"), 10)
	filter := bson.M{
		"deleted_time": nil,
	}

	offset := (page - 1) * pageSize
	res := models.SupplierDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.([]models.Supplier)

	c.JSON(res.Code, res)
}

func UpdateSupplier(c *gin.Context) {
	var req models.Supplier

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

	update := res.Data.(*models.Supplier)

	update.Name = req.Name
	update.Phone = req.Phone
	update.Email = req.Email
	update.Address = req.Address

	updating := models.ProductDB.UpdateOne(filter, update)

	c.JSON(updating.Code, updating)
}

func DeleteSupplier(c *gin.Context) {
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

	res := models.SupplierDB.DeleteOne(filter)

	c.JSON(res.Code, res)
}
