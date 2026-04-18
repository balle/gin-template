package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/balle/gin-template/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("Cannot load .env file")
// 	}
// }

func insertData(db *pgx.Conn) {
	repo := repository.New(db)
	t, _ := time.Parse("2006-01-02 15:04:05", "1993-12-10 00:00:00")

	err := repo.InsertGame(context.Background(), repository.InsertGameParams{Name: "Doom", CreatedDate: pgtype.Timestamptz{Time: t, Valid: true}})

	if err != nil {
		log.Fatalf("Cannot insert game into db: %v", err)
	}

	fmt.Println("Inserted game")
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName))

	if err != nil {
		log.Fatalf("User %s cannot connect to db %s: %v\n", dbUser, dbName, err)
	}

	defer db.Close(context.Background())
	log.Printf("Connected to database %s on %s.", dbName, dbHost)

	//insertData()

	handler := gin.Default()

	handler.GET("/", func(ctx *gin.Context) {
		repo := repository.New(db)
		games, _ := repo.GetAllGames(context.Background())

		ctx.JSON(http.StatusOK, gin.H{
			"message": games,
		})
	})

	handler.Run("0.0.0.0:8000")
}
