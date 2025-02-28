package handlers

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"shop_backend/models"
	"shop_backend/database"
)

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

type UserInfoResponse struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
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
	user, err := database.FindUserByUsername(registerRequest.Username)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	if err != database.ErrUserNotFound {
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
	err = database.CreateUser(registerRequest.Username, hash, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}
	user, err = database.FindUserByUsername(registerRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}
	sessionId := database.CreateSession(user.Id)
	ctx.SetCookie("sessionId", sessionId, 24*3600, "", "", false, true)   // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Login(ctx *gin.Context) {
	var LoginRequest LoginRequest
	if err := ctx.BindJSON(&LoginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	user, err := database.FindUserByUsername(LoginRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(LoginRequest.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}
	sessionId := database.CreateSession(user.Id)
	ctx.SetCookie("sessionId", sessionId, 24*3600, "", "", false, true)   // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Logout(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	database.RemoveSession(sessionId)
	ctx.SetCookie("sessionId", "", 0, "", "", false, true)  // !!!! secure should be true
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func UserInfo(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if (!exists) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	user2 := user.(models.User)
	response := UserInfoResponse{Id: user2.Id, Username: user2.Username, Role: user2.Role}
	ctx.JSON(http.StatusAccepted, response)
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
	err := bcrypt.CompareHashAndPassword([]byte(user.(models.User).PasswordHash), []byte(changePasswordRequest.OldPassword))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	err = database.UpdatePasswordHash(user.(models.User).Id, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
