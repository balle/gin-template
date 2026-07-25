package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/balle/gin-template/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getAllGames(db *gorm.DB) ([]models.Game, error) {
	var games []models.Game
	result := db.Find(&games)

	if result.Error != nil {
		return nil, fmt.Errorf("Could not fetch games: %w", result.Error)
	} else {
		return games, nil
	}
}

func ListAllGames(ctx *gin.Context, db *gorm.DB) {
	games, err := getAllGames(db)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch games",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"games": games,
		})
	}
}

func ViewAllGames(ctx *gin.Context, db *gorm.DB) {
	games, err := getAllGames(db)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch games",
		})
	} else {
		ctx.HTML(http.StatusOK, "game/list.html", gin.H{
			"Games": games,
		})
	}
}

func CreateGame(ctx *gin.Context, db *gorm.DB) {
	input := models.GameInput{}
	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"details": err.Error(),
		})
		return
	}

	db.Find(&input.Game.Gamesystems, input.GamesystemIDs)

	if err := db.Create(&input.Game).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save game", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"error": false, "msg": input.Game})
}

func GetGame(ctx *gin.Context, db *gorm.DB) {
	game := models.Game{}
	inputId := ctx.Param("id")
	id, err := strconv.Atoi(inputId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}

	result := db.Preload("Gamesystems").First(&game, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Game not found",
			"details": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"error": false, "game": game})
}

func UpdateGame(ctx *gin.Context, db *gorm.DB) {
	game := models.Game{}
	inputId := ctx.Param("id")
	id, err := strconv.Atoi(inputId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}

	input := models.GameInput{}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"details": err.Error(),
		})
		return
	}

	result := db.Preload("Gamesystems").First(&game, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Game not found",
			"details": result.Error.Error(),
		})
		return
	}

	game.Name = input.Name
	game.StartedDate = input.StartedDate
	game.FinishedDate = input.FinishedDate
	game.Played = input.Played
	game.Description = input.Description
	game.DownloadOnly = input.DownloadOnly
	game.Rating = input.Rating
	game.ReleaseDate = input.ReleaseDate

	var systems []models.Gamesystem
	if len(input.GamesystemIDs) > 0 {
		if err := db.Find(&systems, input.GamesystemIDs).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid gamesystem IDs",
				"details": err.Error(),
			})
			return
		}
	}

	if err := db.Save(&game).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update game", "details": err.Error()})
		return
	}

	if err := db.Model(&game).Association("Gamesystems").Replace(systems); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update gamesystems", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Updated successfully",
		"game": game,
	})
}

func DeleteGame(ctx *gin.Context, db *gorm.DB) {
	game := models.Game{}
	inputId := ctx.Param("id")
	id, err := strconv.Atoi(inputId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})
		return
	}

	result := db.Preload("Gamesystems").First(&game, id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Game not found",
			"details": result.Error.Error(),
		})
		return
	}

	db.Delete(&game)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Deleted successfully",
	})
}
