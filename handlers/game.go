package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/balle/gin-template/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Internal function to get a game by its id
func getGame(ctx *gin.Context, db *gorm.DB) (*models.Game, error) {
	game := models.Game{}
	result := db.Preload("Gamesystems").First(&game, "id = ?", ctx.Param("id"))

	if result.Error != nil {
		return nil, fmt.Errorf("Failed fetching game by id %d: %v", ctx.Param("id"), result.Error)
	}

	return &game, nil
}

// Internal function to get a list of all games
func getAllGames(db *gorm.DB) ([]models.Game, error) {
	var games []models.Game
	result := db.Find(&games)

	if result.Error != nil {
		return nil, fmt.Errorf("Could not fetch games: %w", result.Error)
	} else {
		return games, nil
	}
}

// ListAllGames godoc
// @Summary      Get a list of all games
// @Description  Get a list of all games
// @Tags         games
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games [get]
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

// SearchGames godoc
// @Summary      Search games by name
// @Description  Search games by name
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        query  name    string  true  "name to search for"  example(doom)
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games/search [get]
func SearchGames(ctx *gin.Context, db *gorm.DB) {
	var games []models.Game
	searchTerm := "%" + strings.ToLower(ctx.Query("name")) + "%"
	err := db.Where("LOWER(name) LIKE ?", searchTerm).Find(&games).Error

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

// Return all games in the html list template
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

// CreateGame godoc
// @Summary      Create a new game
// @Description  Create a new game with given game data
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        game  body      models.GameInput true  "Game data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games/ [post]
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

// GetGame godoc
// @Summary      Get a single game
// @Description  Fetch a single game by its id
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "Game ID"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games/{id} [get]
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

// UpdateGame godoc
// @Summary      Update a game
// @Description  Updates an existing game and its gamesystems
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "Game ID"
// @Param        game  body      models.GameInput true  "Update data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games/{id} [patch]
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

// DeleteGame godoc
// @Summary      Delete a game
// @Description  Delete a game
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "Game ID"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /games/{id} [delete]
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
