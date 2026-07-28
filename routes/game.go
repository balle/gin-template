package routes

import (
	"html/template"
	"net/http"

	"github.com/balle/gin-template/handlers"
	"github.com/balle/gin-template/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountGamesRoutes(handler *gin.Engine, db *gorm.DB, tmpl *template.Template) *gin.Engine {
	handler.SetHTMLTemplate(tmpl)

	handler.Static("/images", "static/images")

	handler.GET("/", func(ctx *gin.Context) {
		handlers.ViewAllGames(ctx, db)
	})

	handler.GET("/games", func(ctx *gin.Context) {
		handlers.ListAllGames(ctx, db)
	})

	authorized := handler.Group("/games")
	//authorized.Use(middleware.APIKeyMiddleware())
	authorized.Use(middleware.JWTAPIKeyMiddleware())
	{
		authorized.GET("/search", func(ctx *gin.Context) {
			handlers.SearchGames(ctx, db)
		})

		authorized.GET("/:id", func(ctx *gin.Context) {
			handlers.GetGame(ctx, db)
		})

		authorized.POST("/", func(ctx *gin.Context) {
			handlers.CreateGame(ctx, db)
		})

		authorized.PUT("/:id", func(ctx *gin.Context) {
			handlers.UpdateGame(ctx, db)
		})

		authorized.PATCH("/:id", func(ctx *gin.Context) {
			handlers.UpdateGame(ctx, db)
		})

		authorized.DELETE("/:id", func(ctx *gin.Context) {
			handlers.DeleteGame(ctx, db)
		})
	}

	handler.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})
	return handler
}
