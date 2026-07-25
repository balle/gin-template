package routes

import (
	"html/template"
	"net/http"

	"github.com/balle/gin-template/routes/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountRoutes(db *gorm.DB, tmpl *template.Template) *gin.Engine {
	handler := gin.Default()
	handler.SetHTMLTemplate(tmpl)

	handler.Static("/images", "static/images")

	handler.GET("/", func(ctx *gin.Context) {
		handlers.ViewAllGames(ctx, db)
	})

	handler.GET("/games", func(ctx *gin.Context) {
		handlers.ListAllGames(ctx, db)
	})

	handler.GET("/games/:id", func(ctx *gin.Context) {
		handlers.GetGame(ctx, db)
	})

	handler.POST("/games", func(ctx *gin.Context) {
		handlers.CreateGame(ctx, db)
	})

	handler.PATCH("/games/:id", func(ctx *gin.Context) {
		handlers.UpdateGame(ctx, db)
	})

	handler.DELETE("/games/:id", func(ctx *gin.Context) {
		handlers.DeleteGame(ctx, db)
	})

	handler.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})
	return handler
}
