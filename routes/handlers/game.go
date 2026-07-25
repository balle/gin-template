package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/balle/gin-template/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getGame(ctx *gin.Context, db *gorm.DB) (*models.Game, error) {
	game := models.Game{}
	result := db.Preload("Gamesystems").First(&game, "id = ?", ctx.Param("id"))

	if result.Error != nil {
		return nil, fmt.Errorf("Failed fetching game by id %d: %v", ctx.Param("id"), result.Error)
	}

	return &game, nil
}

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
			"error": "Validation error",
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
	game, err := getGame(ctx, db)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"error": false, "game": game})
}

func UpdateGame(ctx *gin.Context, db *gorm.DB) {
	input := models.GameInput{}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"details": err.Error(),
		})
		return
	}

	game, err := getGame(ctx, db)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cannot find game",
		})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
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
			if err := tx.Find(&systems, input.GamesystemIDs).Error; err != nil {
				return errors.New("Invalid gamesystem IDs")
			}
		}

		if err := tx.Save(game).Error; err != nil {
			return errors.New("Could not update game")
		}

		if err := tx.Model(game).Association("Gamesystems").Replace(systems); err != nil {
			return errors.New("Could not update gamesystems")
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update game",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "Updated successfully",
			"game": game,
		})
	}
}

func DeleteGame(ctx *gin.Context, db *gorm.DB) {
	game, err := getGame(ctx, db)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Could not find game",
		})
		return
	}

	db.Delete(&game)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Deleted successfully",
	})
}
