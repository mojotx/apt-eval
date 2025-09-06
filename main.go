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

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get data directory - use ./data as default
	dataDir := getEnv("DATA_DIR", filepath.Join(".", "data"))

	// Setup database
	database, err := db.New(dataDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.Close()

	// Setup router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Root route to serve the main HTML page
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
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

	// Setup and start server
	httpsPort := getEnv("PORT", "8443")
	httpPort := getEnv("HTTP_PORT", "8080")
	certFile := getEnv("CERT_FILE", "./certs/wildcard.crt")
	keyFile := getEnv("KEY_FILE", "./certs/wildcard.key")

	// Configure TLS settings for HTTPS server
	srv := &http.Server{
		Addr:      ":" + httpsPort,
		Handler:   router,
		TLSConfig: getTLSConfig(),
	}

	// Setup HTTP server to redirect to HTTPS
	httpSrv := &http.Server{
		Addr: ":" + httpPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.Host, ":")[0]
			target := "https://" + host + ":" + httpsPort + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}),
	}

	// Run HTTP server in a goroutine for redirects
	go func() {
		log.Info().Str("port", httpPort).Msg("Starting HTTP server (for redirects)")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
		}
	}()

	// Run HTTPS server in a goroutine
	go func() {
		log.Info().Str("port", httpsPort).Msg("Starting secure server (HTTPS)")
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start secure server")
		}
	}()

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
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTPS server forced to shutdown")
	}

	// Shutdown HTTP server
	log.Info().Msg("Shutting down HTTP server...")
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server forced to shutdown")
	}

	log.Info().Msg("Servers exited properly")
}

// getEnv gets an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getTLSConfig returns a TLS configuration with modern security settings
func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
