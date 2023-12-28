package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")

		if bearerToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"});
			ctx.Abort()
			return
		}
	

		bearerSplit := strings.Split(bearerToken, " ")

		tokenLen := len(bearerSplit)
		token := bearerSplit[1]
		name := bearerSplit[0]

		if tokenLen != 2 || strings.ToLower(name) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access [Invalid token type provided]"});
			ctx.Abort()
			return
		}

		userId, err := tokenController.VerifyToken(token);
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()});
			ctx.Abort();
			return;
		}

		ctx.Set("user_id", userId)
	}
}