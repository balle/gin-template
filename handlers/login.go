package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/balle/gin-template/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Login(ctx *gin.Context, db *gorm.DB) {
	input := models.User{}
	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"error":   err.Error(),
		})
		return
	}

	result := db.Where("username = ? and password = ?", input.Username, input.Password).First(&input)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Login failed",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Login failed",
		})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		Username: input.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-template",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil || tokenString == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"messgae": "Cannot create JWT token",
			"error":   err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": expirationTime,
	})
}
