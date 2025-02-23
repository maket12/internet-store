package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
  _ "modernc.org/sqlite"
)

type Product struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Description string `json:"description"`
  Image string `json:"image"`
  Price float64 `json:"price"`
  Available int `json:"available"`
}

type AddProductRequest struct {
  Name string `json:"name"`
  Description string `json:"description"`
  Image string `json:"image"`
  Price float64 `json:"price"`
  Available int `json:"available"`
}

type UpdateProductRequest struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Description string `json:"description"`
  Image string `json:"image"`
  Price float64 `json:"price"`
  Available int `json:"available"`
}

type RemoveProductRequest struct {
  Id string `json:"id"`
}

var database *sql.DB

func InitDatabase() {
  database.Exec(`
  CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    image TEXT,
    price REAL,
    available INTEGER
  )
  `)
}

func GetProducts(ctx *gin.Context) {
  rows, err := database.Query("SELECT * FROM products")
  if err != nil {
    fmt.Print(err)
    fmt.Print("failed to query")
    ctx.JSON(500, gin.H{"error": "failed to query"})
    return
  }
  products := []Product{}
  for rows.Next() {
    var p Product
    rows.Scan(&p.Id, &p.Name, &p.Description, &p.Image, &p.Price, &p.Available)
    products = append(products, p)
  }
  ctx.JSON(200, products)
}

func AddProduct(ctx *gin.Context) {
  var addProductRequest AddProductRequest
  if err := ctx.BindJSON(&addProductRequest); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
    return
  }
  id := uuid.New().String()
  _, err := database.Exec("INSERT INTO products (id, name, description, image, price, available) VALUES (?, ?, ?, ?, ?, ?)", id, addProductRequest.Name, addProductRequest.Description, addProductRequest.Image, addProductRequest.Price, addProductRequest.Available)
  if err != nil {
    fmt.Print("failed to insert")
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
    return
  }
  ctx.JSON(http.StatusOK, gin.H{"message": "success", "id": id})
}

func UpdateProduct(ctx *gin.Context) {
  var updateProductRequest UpdateProductRequest
  if err := ctx.BindJSON(&updateProductRequest); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
    return
  }
  _, err := database.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, available = ? WHERE id = ?", updateProductRequest.Name, updateProductRequest.Description, updateProductRequest.Image, updateProductRequest.Price, updateProductRequest.Available, updateProductRequest.Id)
  if err != nil {
    fmt.Print("failed to update")
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
  _, err := database.Exec("DELETE FROM products WHERE id = ?", removeProductRequest.Id)
  if err != nil {
    fmt.Print("failed to delete")
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
    return
  }
  ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func main() {
  var err error
  database, err = sql.Open("sqlite", "store.db")
  if err != nil {
    fmt.Print("failed to open db")
    return
  }
  InitDatabase()
  r := gin.Default()
  r.StaticFile("/test", "./test.html")
  r.GET("/api/v1/get_products", GetProducts)
  r.POST("/api/v1/add_product", AddProduct)
  r.POST("/api/v1/update_product", UpdateProduct)
  r.POST("/api/v1/remove_product", RemoveProduct)
  r.Run() // listen and serve on 0.0.0.0:8080
}
