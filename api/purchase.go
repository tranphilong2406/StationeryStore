package api

import (
	"StoreServer/models"
	"StoreServer/utils"
	myerror "StoreServer/utils/error"
	"StoreServer/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreatePurchase(c *gin.Context) {
	var req models.Purchase
	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	fmt.Println("CreatePurchase request:", req)

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	res := models.PurchaseDB.Create(req)
	c.JSON(res.Code, res)
}

func ApprovePurchase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.MyResponse.Error(c, myerror.EmptyParam())
		return
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	purchase := models.PurchaseDB.QueryOne(filter)
	if purchase.Code != http.StatusOK {
		c.JSON(purchase.Code, purchase)
		return
	}

	purchaseData := purchase.Data.(*models.Purchase)
	purchaseData.Status = models.Approved

	for idx := range len(purchaseData.Products) {
		filter := bson.M{
			"_id": purchaseData.Products[idx].ID,
		}

		prod := models.ProductDB.QueryOne(filter)
		if prod.Code != http.StatusOK {
			c.JSON(prod.Code, prod)
			return
		}

		update := prod.Data.(*models.Product)

		update.Stock += purchaseData.Products[idx].Quantity
		update.BuyPrice = purchaseData.Products[idx].Buy

		updated := models.ProductDB.UpdateOne(filter, update)
		if updated.Code != http.StatusOK {
			c.JSON(updated.Code, updated)
			return
		}
	}

	res := models.PurchaseDB.UpdateOne(filter, purchaseData)
	c.JSON(res.Code, res)
}

func RejectPurchase(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.MyResponse.Error(c, myerror.EmptyParam())
		return
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	purchase := models.PurchaseDB.QueryOne(filter)
	if purchase.Code != http.StatusOK {
		c.JSON(purchase.Code, purchase)
		return
	}

	purchaseData := purchase.Data.(*models.Purchase)
	purchaseData.Status = models.Rejected

	res := models.PurchaseDB.UpdateOne(filter, purchaseData)
	c.JSON(res.Code, res)
}

func GetPurchase(c *gin.Context) {
	page := utils.ParseInt(c.Query("page"), 1)
	pageSize := utils.ParseInt(c.Query("page_size"), 10)

	filter := bson.M{
		"deleted_time": nil,
	}

	offset := (page - 1) * pageSize
	res := models.PurchaseDB.Query(filter, offset, pageSize)
	if res.Code != http.StatusOK {
		c.JSON(res.Code, res)
		return
	}

	res.Data = res.Data.([]models.Purchase)
	res.Page = page
	res.PageSize = pageSize

	c.JSON(res.Code, res)
}

func GetPurchaseByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.MyResponse.Error(c, myerror.EmptyParam())
		return
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusInternalServerError, err))
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	result := models.PurchaseDB.QueryOne(filter)

	c.JSON(result.Code, result)
}

func UpdatePurchase(c *gin.Context) {
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

	var req models.Purchase
	if err := c.ShouldBind(&req); err != nil {
		response.MyResponse.Error(c, myerror.AnyError(http.StatusBadRequest, err))
		return
	}

	req.ID = objID

	if ok := req.Validate(); ok.Code != http.StatusOK {
		c.JSON(ok.Code, ok)
		return
	}

	filter := bson.M{
		"_id": objID,
	}

	res := models.PurchaseDB.UpdateOne(filter, req)
	c.JSON(res.Code, res)
}

func DeletePurchase(c *gin.Context) {
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

	res := models.PurchaseDB.DeleteOne(filter)

	c.JSON(res.Code, res)
}
