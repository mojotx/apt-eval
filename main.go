package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojotx/apt-eval/db"
	"github.com/mojotx/apt-eval/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// App holds the application components
type App struct {
	DB       *db.DB
	Router   *gin.Engine
	HTTPSrv  *http.Server
	RedirSrv *http.Server
	Config   AppConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	DataDir    string
	HTTPPort   string
	HTTPSPort  string
	CertFile   string
	KeyFile    string
	StaticPath string
}

func main() {
	setupLogging()

	// Initialize application config
	config := loadConfig()

	// Create and initialize the app
	app, err := initApp(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}
	defer app.DB.Close()

	// Start the servers
	startServers(app)

	// Wait for shutdown signal and handle graceful shutdown
	handleShutdown(app)
}

// setupLogging configures the application logging
func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// loadConfig loads application configuration from environment variables
func loadConfig() AppConfig {
	return AppConfig{
		DataDir:    getEnv("DATA_DIR", filepath.Join(".", "data")),
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
		HTTPSPort:  getEnv("PORT", "8443"),
		CertFile:   getEnv("CERT_FILE", "./certs/wildcard.crt"),
		KeyFile:    getEnv("KEY_FILE", "./certs/wildcard.key"),
		StaticPath: "./static",
	}
}

// initApp initializes the application components
func initApp(config AppConfig) (*App, error) {
	// Initialize database
	database, err := db.New(config.DataDir)
	if err != nil {
		return nil, err
	}

	// Setup router with routes
	router := setupRouter(database, config)

	// Create app instance
	app := &App{
		DB:     database,
		Router: router,
		Config: config,
	}

	// Configure HTTP and HTTPS servers
	setupServers(app)

	return app, nil
}

// setupRouter configures the Gin router with all routes
func setupRouter(database *db.DB, config AppConfig) *gin.Engine {
	router := gin.Default()

	// Serve static files
	router.Static("/static", config.StaticPath)

	// Root route to serve the main HTML page
	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(config.StaticPath, "index.html"))
	})

	// Setup API routes
	apartmentHandler := handlers.NewApartmentHandler(database)
	apartmentHandler.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"time":   time.Now().Unix(),
		})
	})

	return router
}

// setupServers configures the HTTP and HTTPS servers
func setupServers(app *App) {
	// Configure TLS settings for HTTPS server
	app.HTTPSrv = &http.Server{
		Addr:      ":" + app.Config.HTTPSPort,
		Handler:   app.Router,
		TLSConfig: getTLSConfig(),
	}

	// Setup HTTP server to redirect to HTTPS
	app.RedirSrv = &http.Server{
		Addr: ":" + app.Config.HTTPPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.Host, ":")[0]
			target := "https://" + host + ":" + app.Config.HTTPSPort + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}),
	}
}

// startServers starts both HTTP and HTTPS servers
func startServers(app *App) {
	// Run HTTP server in a goroutine for redirects
	go func() {
		log.Info().Str("port", app.Config.HTTPPort).Msg("Starting HTTP server (for redirects)")
		if err := app.RedirSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	// Run HTTPS server in a goroutine
	go func() {
		log.Info().Str("port", app.Config.HTTPSPort).Msg("Starting secure server (HTTPS)")
		if err := app.HTTPSrv.ListenAndServeTLS(app.Config.CertFile, app.Config.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start secure server")
		}
	}()
}

// handleShutdown waits for termination signal and performs graceful shutdown
func handleShutdown(app *App) {
	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down servers...")

	// Give servers 5 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTPS server
	log.Info().Msg("Shutting down HTTPS server...")
	if err := app.HTTPSrv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTPS server forced to shutdown")
	}

	// Shutdown HTTP server
	log.Info().Msg("Shutting down HTTP server...")
	if err := app.RedirSrv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server forced to shutdown")
	}

	log.Info().Msg("Servers exited properly")
}

// getEnv returns environment variable value or fallback if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getTLSConfig returns TLS configuration with secure defaults
func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
			tls.X25519,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
