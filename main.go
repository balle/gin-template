package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/balle/gin-template/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", dbHost, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database")
	}

	db.AutoMigrate(&models.Game{})

	//defer db.Close()
	log.Printf("Connected to database %s on %s.", dbName, dbHost)

	handler := gin.Default()

	handler.GET("/", func(ctx *gin.Context) {
		var games []models.Game
		result := db.Find(&games)

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not fetch games",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": games,
			})
		}
	})

	handler.Run("0.0.0.0:8000")
}
