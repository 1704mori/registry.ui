package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/1704mori/registry.ui/internal/api"
	"github.com/1704mori/registry.ui/internal/config"
	"github.com/1704mori/registry.ui/internal/handlers"
)

func main() {
	// Load environment variables from .env file if it exists
	godotenv.Load()

	// Initialize configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize registry client
	registryClient, err := api.NewClient(cfg.RegistryURL, cfg.RegistryUsername, cfg.RegistryPassword, cfg.RegistryInsecure)
	if err != nil {
		log.Fatalf("Failed to initialize registry client: %v", err)
	}

	// Initialize Echo framework
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Static files
	e.Static("/static", "static")

	// Initialize handlers
	h := handlers.NewHandlers(registryClient, cfg)

	// Routes
	e.GET("/", h.Dashboard)
	e.GET("/images", h.ListImages)
	// Register a single catch-all route for GET and DELETE requests.
	e.GET("/images/*", h.GetImage)
	e.GET("/image-tags/*", h.HtmxListTags)
	e.GET("/image-tag/:tag/*", h.GetTag)
	e.DELETE("/image-tag/:tag/*", h.DeleteTag)
	e.GET("/settings", h.Settings)
	e.POST("/settings", h.UpdateSettings)

	// HTMX specific endpoints
	e.GET("/htmx/images", h.HtmxListImages)
	e.GET("/htmx/image-tags/*", h.HtmxListTags)
	e.GET("/htmx/theme-toggle", h.HtmxThemeToggle)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	s := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on %s", serverAddr)
	if err := e.StartServer(s); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
