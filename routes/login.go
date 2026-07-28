package routes

import (
	"html/template"

	"github.com/balle/gin-template/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MountLoginRoutes(handler *gin.Engine, db *gorm.DB, tmpl *template.Template) *gin.Engine {
	handler.GET("/login", func(ctx *gin.Context) {
		handlers.Login(ctx, db)
	})

	return handler
}
