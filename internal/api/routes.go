package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skb1129/calculator-app/internal/api/handlers"
	"github.com/skb1129/calculator-app/internal/api/middleware"
)

// SetupRouter configures the Gin engine with routes and middleware.
func SetupRouter(staticDir string) *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORS())

	// API routes group
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", handlers.HealthCheck)
		apiV1.POST("/calculate", handlers.Calculate)
	}

	// Serve static files if directory is provided and exists
	if staticDir != "" {
		if fi, err := os.Stat(staticDir); err == nil && fi.IsDir() {
			// Serve assets
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path

				// Don't intercept API calls
				if strings.HasPrefix(path, "/api") {
					c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
					return
				}

				// Check if the requested file exists
				targetFile := filepath.Join(staticDir, filepath.Clean(path))
				if fileInfo, err := os.Stat(targetFile); err == nil && !fileInfo.IsDir() {
					c.File(targetFile)
					return
				}

				// Otherwise fallback to index.html for SPA client-side routing
				indexFile := filepath.Join(staticDir, "index.html")
				c.File(indexFile)
			})
		}
	}

	return r
}
