package api

import (
	"StoreServer/models"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var req models.User
	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	req.Password = string(passHash)

	res := models.AuthDB.Create(req)
	c.JSON(res.Code, res)
}

func UpdateUser(c *gin.Context) {
	var req models.User

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
		"_id": objID,
	}

	res := models.AuthDB.UpdateOne(filter, req)
	c.JSON(res.Code, res)
}
