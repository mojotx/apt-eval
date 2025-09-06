package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mojotx/apt-eval/db"
	"github.com/mojotx/apt-eval/models"
	"github.com/rs/zerolog/log"
)

// ApartmentHandler handles apartment-related requests
type ApartmentHandler struct {
	db *db.DB
}

// NewApartmentHandler creates a new apartment handler
func NewApartmentHandler(db *db.DB) *ApartmentHandler {
	return &ApartmentHandler{
		db: db,
	}
}

// Create handles the creation of a new apartment evaluation
func (h *ApartmentHandler) Create(c *gin.Context) {
	var request models.ApartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apartment, err := h.db.CreateApartment(&request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create apartment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create apartment"})
		return
	}

	c.JSON(http.StatusCreated, apartment)
}

// Get handles retrieving an apartment by ID
func (h *ApartmentHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Invalid apartment ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid apartment ID"})
		return
	}

	apartment, err := h.db.GetApartment(id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to get apartment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get apartment"})
		return
	}

	if apartment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
		return
	}

	c.JSON(http.StatusOK, apartment)
}

// List handles retrieving all apartments
func (h *ApartmentHandler) List(c *gin.Context) {
	apartments, err := h.db.ListApartments()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list apartments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list apartments"})
		return
	}

	c.JSON(http.StatusOK, apartments)
}

// Update handles updating an apartment
func (h *ApartmentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Invalid apartment ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid apartment ID"})
		return
	}

	var request models.ApartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apartment, err := h.db.UpdateApartment(id, &request)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to update apartment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update apartment"})
		return
	}

	if apartment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
		return
	}

	c.JSON(http.StatusOK, apartment)
}

// Delete handles deleting an apartment
func (h *ApartmentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("Invalid apartment ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid apartment ID"})
		return
	}

	err = h.db.DeleteApartment(id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to delete apartment")
		if err.Error() == "apartment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Apartment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete apartment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RegisterRoutes registers all apartment-related routes
func (h *ApartmentHandler) RegisterRoutes(router *gin.Engine) {
	apartments := router.Group("/api/apartments")
	{
		apartments.POST("", h.Create)
		apartments.GET("", h.List)
		apartments.GET("/:id", h.Get)
		apartments.PUT("/:id", h.Update)
		apartments.DELETE("/:id", h.Delete)
	}
}
