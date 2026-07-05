package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"url-shortener/models"
	"url-shortener/utils"
)

// URLHandler bundles together everything the handlers need.
// Right now that's just the in-memory store, but this pattern makes
// it easy to add things like a logger or config later.
type URLHandler struct {
	Store *models.URLStore
}

// NewURLHandler creates a new URLHandler wired up with the given store.
func NewURLHandler(store *models.URLStore) *URLHandler {
	return &URLHandler{Store: store}
}

// ShortenRequest is the expected JSON body for POST /shorten.
// `binding:"required"` tells Gin to reject the request if "url" is missing.
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// isValidURL checks that a string is a well-formed http/https URL.
func isValidURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

// ShortenURL handles POST /shorten.
// It validates the input, generates a unique short code, saves the
// mapping, and returns the full short URL to the client.
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req ShortenRequest

	// ShouldBindJSON parses the request body into req and validates it
	// against the struct tags (e.g. "required").
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body. Expected JSON: {\"url\": \"https://example.com\"}",
		})
		return
	}

	if !isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL. It must start with http:// or https://",
		})
		return
	}

	// Generate a short code, retrying in the rare case of a collision.
	code := utils.GenerateShortCode()
	for h.Store.Exists(code) {
		code = utils.GenerateShortCode()
	}

	h.Store.Save(code, req.URL)

	shortURL := "http://localhost:8080/" + code
	c.JSON(http.StatusCreated, gin.H{"short_url": shortURL})
}

// RedirectURL handles GET /:code.
// It looks up the original URL, increments the click counter, and
// issues an HTTP 302 redirect to the browser/client.
func (h *URLHandler) RedirectURL(c *gin.Context) {
	code := c.Param("code")

	entry, found := h.Store.Get(code)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	h.Store.IncrementClicks(code)

	// http.StatusFound = 302 (temporary redirect)
	c.Redirect(http.StatusFound, entry.OriginalURL)
}

// ListURLs handles GET /urls.
// Returns every stored URL along with its short code and click count.
func (h *URLHandler) ListURLs(c *gin.Context) {
	urls := h.Store.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"count": len(urls),
		"urls":  urls,
	})
}

// DeleteURL handles DELETE /urls/:code.
// Removes a short URL entry if it exists.
func (h *URLHandler) DeleteURL(c *gin.Context) {
	code := c.Param("code")

	if !h.Store.Delete(code) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Short URL deleted successfully"})
}
