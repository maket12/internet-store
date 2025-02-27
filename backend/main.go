package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Available   int     `json:"available"`
}

type AddProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Available   int     `json:"available"`
}

type UpdateProductRequest struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Available   int     `json:"available"`
}

type RemoveProductRequest struct {
	Id string `json:"id"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
}

type UserWithPassword struct {
	User
	PasswordHash	string
}

type Session struct {
	SessionId  string `json:"sessionId"`
	UserId     int    `json:"userId"`
	ExpireTime int64  `json:"expireTime"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var ErrInvalidSession = errors.New("invalid session")
var ErrSessionExpired = errors.New("session expired")
var ErrInvalidUser = errors.New("invalid user")

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
	database.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		passwordHash TEXT,
		role TEXT
	)
	`)
	database.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		sessionId TEXT PRIMARY KEY,
		userId INTEGER,
		expireTime INTEGER,
		FOREIGN KEY(userId) REFERENCES users(id)
	)
	`)
}

func CreateSession(userId int) string {
	sessionId := uuid.New().String()
	expireTime := time.Now().Add(time.Hour * 24).Unix()
	database.Exec("INSERT INTO sessions (sessionId, userId, expireTime) VALUES (?, ?, ?)", sessionId, userId, expireTime)
	return sessionId
}

func GetUserBySession(sessionId string) (User, error) {
	var session Session
	session.SessionId = sessionId
	err := database.QueryRow("SELECT userId, expireTime FROM sessions WHERE sessionId = ?", sessionId).Scan(&session.UserId, &session.ExpireTime)
	if err != nil {
		return User{}, ErrInvalidSession
	}
	if session.ExpireTime < time.Now().Unix() {
		database.Exec("DELETE FROM sessions WHERE sessionId = ?", sessionId)
		return User{}, ErrSessionExpired
	}
	var user User
	err = database.QueryRow("SELECT id, username, role FROM users WHERE id = ?", session.UserId).Scan(&user.Id, &user.Username, &user.Role)
	if err != nil {
		return User{}, ErrInvalidUser
	}
	return user, nil
}

func Register(ctx *gin.Context) {
	var registerRequest RegisterRequest
	if err := ctx.BindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	registerRequest.Username = strings.TrimSpace(registerRequest.Username)
	if len(registerRequest.Username) < 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username too short"})
		return
	}
	registerRequest.Password = strings.TrimSpace(registerRequest.Password)
	if len(registerRequest.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		return
	}
	var id int
	err := database.QueryRow("SELECT id FROM users WHERE username = ?", registerRequest.Username).Scan(&id)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	if err != sql.ErrNoRows {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query"})
		return
	}
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	role := "user"
	if registerRequest.Username == "admin" {
		role = "admin"
	}
	_, err = database.Exec("INSERT INTO users (username, passwordHash, role) VALUES (?, ?, ?)", registerRequest.Username, hash, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}
	var user User
	err = database.QueryRow("SELECT id, username, role FROM users WHERE username = ?", registerRequest.Username).Scan(&user.Id, &user.Username, &user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}
	sessionId := CreateSession(user.Id)
	ctx.SetCookie("sessionId", sessionId, 24*3600, "", "", false, true)   // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Login(ctx *gin.Context) {
	var LoginRequest LoginRequest
	if err := ctx.BindJSON(&LoginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var user UserWithPassword
	err := database.QueryRow("SELECT id, username, passwordHash, role FROM users WHERE username = ?", LoginRequest.Username).Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(LoginRequest.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}
	sessionId := CreateSession(user.Id)
	ctx.SetCookie("sessionId", sessionId, 24*3600, "", "", false, true)   // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Logout(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	database.Exec("DELETE FROM sessions WHERE sessionId = ?", sessionId)
	ctx.SetCookie("sessionId", "", 0, "", "", false, true)  // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func AuthMiddleware(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	user, err := GetUserBySession(sessionId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}

func UserInfo(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if (!exists) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

func ChangePassword(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if (!exists) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	var changePasswordRequest ChangePasswordRequest
	if err := ctx.BindJSON(&changePasswordRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	changePasswordRequest.NewPassword = strings.TrimSpace(changePasswordRequest.NewPassword)
	if len(changePasswordRequest.NewPassword) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.(UserWithPassword).PasswordHash), []byte(changePasswordRequest.OldPassword))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
	database.Exec("UPDATE passwordHash FROM users WHERE id = ?", hash, user.(User).Id)
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func AdminMiddleware(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if (!exists) {
		return
	}
	if (user.(User).Role != "admin") {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not admin"})
		return
	}
	ctx.Next()
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

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost", "http://150.241.94.206"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
	}))
	engine.GET("/api/v1/get_products", GetProducts)
	engine.POST("/api/v1/register", Register)
	engine.POST("/api/v1/login", Login)

	authGroup := engine.Group("")
	authGroup.Use(AuthMiddleware)
	authGroup.POST("/api/v1/logout", Logout)
	authGroup.POST("/api/v1/user_info", UserInfo)
	authGroup.POST("/api/v1/change_password", ChangePassword)

	adminGroup := authGroup.Group("")
	adminGroup.Use(AdminMiddleware)
	adminGroup.POST("/api/v1/add_product", AddProduct)
	adminGroup.POST("/api/v1/update_product", UpdateProduct)
	adminGroup.POST("/api/v1/remove_product", RemoveProduct)

	engine.Run() // listen and serve on 0.0.0.0:8080
}
