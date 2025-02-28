package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/models"
)

func AdminMiddleware(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if (!exists) {
		return
	}
	if (user.(models.User).Role != "admin") {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not admin"})
		return
	}
	ctx.Next()
}
