package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"shop_backend/models"
	"shop_backend/database"
)

type AddProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Available   int     `json:"available"`
}

type RemoveProductRequest struct {
	Id string `json:"id"`
}

func GetProducts(ctx *gin.Context) {
	products, err := database.GetAllProducts()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to query"})
		return
	}
	ctx.JSON(200, products)
}

func AddProduct(ctx *gin.Context) {
	var addProductRequest AddProductRequest
	if err := ctx.BindJSON(&addProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	p := models.Product{Name: addProductRequest.Name, Description: addProductRequest.Description, Image: addProductRequest.Image, Price: addProductRequest.Price, Available: addProductRequest.Available}
	id, err := database.AddProduct(p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "id": id})
}

func UpdateProduct(ctx *gin.Context) {
	var updateProductRequest models.Product
	if err := ctx.BindJSON(&updateProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	err := database.UpdateProduct(updateProductRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func RemoveProduct(ctx *gin.Context) {
	var removeProductRequest RemoveProductRequest
	if err := ctx.BindJSON(&removeProductRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	err := database.RemoveProduct(removeProductRequest.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
