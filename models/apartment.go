package models

import (
	"time"
)

// Apartment represents an apartment evaluation record
type Apartment struct {
	ID       int64     `json:"id"`
	Address  string    `json:"address" binding:"required"`
	VisitDate time.Time `json:"visit_date"`
	Notes    string    `json:"notes"`
	Rating   int       `json:"rating"` // Rating from 1-5
	Price    float64   `json:"price"`  // Monthly rent/price
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ApartmentRequest is used for creating/updating an apartment record
type ApartmentRequest struct {
	Address  string    `json:"address" binding:"required"`
	VisitDate time.Time `json:"visit_date"`
	Notes    string    `json:"notes"`
	Rating   int       `json:"rating"`
	Price    float64   `json:"price"`
}
