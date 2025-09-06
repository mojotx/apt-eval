package models

import (
	"strings"
	"time"
)

// Apartment represents an apartment evaluation record
type Apartment struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address" binding:"required"`
	VisitDate time.Time `json:"visit_date"`
	Notes     string    `json:"notes"`
	Rating    int       `json:"rating"` // Rating from 1-5
	Price     float64   `json:"price"`  // Monthly rent/price
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CustomTime is a wrapper around time.Time to handle various date formats
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for CustomTime
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return nil
	}

	// Try various time formats
	formats := []string{
		time.RFC3339,          // Standard format: 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05", // No timezone
		"2006-01-02T15:04",    // No seconds or timezone
		"2006-01-02",          // Just date
	}

	var err error
	for _, format := range formats {
		t, parseErr := time.Parse(format, s)
		if parseErr == nil {
			ct.Time = t
			return nil
		}
		err = parseErr
	}

	return err
}

// ApartmentRequest is used for creating/updating an apartment record
type ApartmentRequest struct {
	Address   string     `json:"address" binding:"required"`
	VisitDate CustomTime `json:"visit_date"`
	Notes     string     `json:"notes"`
	Rating    int        `json:"rating"`
	Price     float64    `json:"price"`
}
