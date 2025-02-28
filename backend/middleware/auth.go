package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/database"
)

func AuthMiddleware(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	user, err := database.GetUserBySession(sessionId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}
