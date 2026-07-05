package routes

import (
	"github.com/gin-gonic/gin"

	"url-shortener/handlers"
	"url-shortener/models"
)

// SetupRouter builds and returns a fully configured Gin engine.
// Keeping this separate from main.go makes the app easier to test.
func SetupRouter() *gin.Engine {
	store := models.NewURLStore()
	h := handlers.NewURLHandler(store)

	r := gin.Default() // includes logging + panic-recovery middleware

	r.POST("/shorten", h.ShortenURL)
	r.GET("/urls", h.ListURLs)
	r.DELETE("/urls/:code", h.DeleteURL)

	// Wildcard route — keep this last so it doesn't shadow /shorten or /urls.
	r.GET("/:code", h.RedirectURL)

	return r
}
