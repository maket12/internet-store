package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"shop_backend/database"
	"shop_backend/handlers"
	"shop_backend/middleware"
)

func main() {
	err := database.OpenDatabase()
	if err != nil {
		fmt.Println("failed to open database")
		fmt.Println(err)
		return
	}
	defer database.CloseDatabase()

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://150.241.94.206"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
	}))
	api := engine.Group("/api/v1")
	api.GET("/get_products", handlers.GetProducts)
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	authGroup := api.Group("")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.POST("/logout", handlers.Logout)
	authGroup.POST("/user_info", handlers.UserInfo)
	authGroup.POST("/change_password", handlers.ChangePassword)

	adminGroup := authGroup.Group("")
	adminGroup.Use(middleware.AdminMiddleware)
	adminGroup.POST("/add_product", handlers.AddProduct)
	adminGroup.POST("/update_product", handlers.UpdateProduct)
	adminGroup.POST("/remove_product", handlers.RemoveProduct)

	engine.Run() // listen and serve on 0.0.0.0:8080
}
