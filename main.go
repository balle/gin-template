package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/balle/gin-template/models"
	"github.com/balle/gin-template/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/balle/gin-template/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func insertTestData(db *gorm.DB) {
	t, _ := time.Parse("2006-01-02 15:04:05", "1993-12-10 00:00:00")

	gamesystem := models.Gamesystem{
		Name: "PC",
	}
	result := db.Create(&gamesystem)

	if result.Error != nil {
		log.Fatalf("Cannot insert gamesystem test data: %v", result.Error)
	}

	game := models.Game{
		Name:        "Doom",
		Played:      true,
		Description: "Best game ever",
		Rating:      100,
		ReleaseDate: t,
		Gamesystems: []models.Gamesystem{gamesystem},
	}
	result = db.Create(&game)

	if result.Error != nil {
		log.Fatalf("Cannot insert game test data: %v", result.Error)
	}

	user := models.User{
		Username: "gin",
		Password: "tonic",
	}

	result = db.Create(&user)

	if result.Error != nil {
		log.Fatalf("Cannot insert user test data: %v", result.Error)
	}

	log.Println("Inserted test data.")
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", dbHost, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.Gamesystem{}, &models.Game{}, &models.User{})

	if err != nil {
		log.Fatalf("Db migration failed: %v", err)
	}

	log.Printf("Connected to database %s on %s.", dbName, dbHost)

	// insertTestData(db)

	tmpl, err := template.ParseGlob("templates/game/*.html")

	if err != nil {
		log.Fatalf("Template errors:\n%v\n", err)
	}

	handler := gin.Default()
	handler = routes.MountLoginRoutes(handler, db, tmpl)
	handler = routes.MountGamesRoutes(handler, db, tmpl)
	handler.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.Run("0.0.0.0:8000")
}
