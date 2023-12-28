package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Currencies = map[string]string {
	"USD": "USD",
	"NGN": "NGN",
	"GBP": "GBP",
}

func IsValidCurrency(currency string) bool {
	if _, ok := Currencies[currency]; !ok {
		return false
	}

	return true
}

func GetActiveUser(ctx *gin.Context) (int64, error) {
	userIdFromContext, exist := ctx.Get("user_id");
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return 0, fmt.Errorf("unathorized access");
	}

	userID, ok := userIdFromContext.(int64);
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! something went wrong"})
		return  0, fmt.Errorf("oops! something went wrong");
	}

	return userID, nil
}