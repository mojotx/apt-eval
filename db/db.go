package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mojotx/apt-eval/models"
	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dataDir string) (*DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "apartments.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize database schema
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &DB{db}, nil
}

//go:embed create.sql
var createTableQuery string

// initSchema creates the necessary tables if they don't exist
func initSchema(db *sql.DB) error {

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Info().Msg("Database schema initialized")
	return nil
}

//go:embed insert.sql
var insertApartmentQuery string

// CreateApartment inserts a new apartment record
func (db *DB) CreateApartment(apt *models.ApartmentRequest) (*models.Apartment, error) {
	var apartment models.Apartment
	err := db.QueryRow(
		insertApartmentQuery,
		apt.Address,
		apt.VisitDate.Time,
		apt.Notes,
		apt.Rating,
		apt.Price,
		apt.Floor,
		apt.IsGated,
		apt.HasGarage,
		apt.HasLaundry,
	).Scan(
		&apartment.ID,
		&apartment.Address,
		&apartment.VisitDate,
		&apartment.Notes,
		&apartment.Rating,
		&apartment.Price,
		&apartment.Floor,
		&apartment.IsGated,
		&apartment.HasGarage,
		&apartment.HasLaundry,
		&apartment.CreatedAt,
		&apartment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create apartment: %w", err)
	}

	return &apartment, nil
}

//go:embed get.sql
var getApartmentQuery string

// GetApartment retrieves an apartment by ID
func (db *DB) GetApartment(id int64) (*models.Apartment, error) {

	var apartment models.Apartment
	err := db.QueryRow(getApartmentQuery, id).Scan(
		&apartment.ID,
		&apartment.Address,
		&apartment.VisitDate,
		&apartment.Notes,
		&apartment.Rating,
		&apartment.Price,
		&apartment.Floor,
		&apartment.IsGated,
		&apartment.HasGarage,
		&apartment.HasLaundry,
		&apartment.CreatedAt,
		&apartment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get apartment: %w", err)
	}

	return &apartment, nil
}

//go:embed list.sql
var listApartmentsQuery string

// ListApartments retrieves all apartments
func (db *DB) ListApartments() ([]models.Apartment, error) {

	rows, err := db.Query(listApartmentsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list apartments: %w", err)
	}
	defer rows.Close()

	apartments := []models.Apartment{}
	for rows.Next() {
		var apt models.Apartment
		if err := rows.Scan(
			&apt.ID,
			&apt.Address,
			&apt.VisitDate,
			&apt.Notes,
			&apt.Rating,
			&apt.Price,
			&apt.Floor,
			&apt.IsGated,
			&apt.HasGarage,
			&apt.HasLaundry,
			&apt.CreatedAt,
			&apt.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan apartment row: %w", err)
		}
		apartments = append(apartments, apt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return apartments, nil
}

// UpdateApartment updates an existing apartment
func (db *DB) UpdateApartment(id int64, apt *models.ApartmentRequest) (*models.Apartment, error) {
	query := `
		UPDATE apartments
		SET address = ?, visit_date = ?, notes = ?, rating = ?, price = ?, 
		    floor = ?, is_gated = ?, has_garage = ?, has_laundry = ?, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, address, visit_date, notes, rating, price, floor, is_gated, has_garage, has_laundry, created_at, updated_at
	`

	var apartment models.Apartment
	err := db.QueryRow(
		query,
		apt.Address,
		apt.VisitDate.Time,
		apt.Notes,
		apt.Rating,
		apt.Price,
		apt.Floor,
		apt.IsGated,
		apt.HasGarage,
		apt.HasLaundry,
		id,
	).Scan(
		&apartment.ID,
		&apartment.Address,
		&apartment.VisitDate,
		&apartment.Notes,
		&apartment.Rating,
		&apartment.Price,
		&apartment.Floor,
		&apartment.IsGated,
		&apartment.HasGarage,
		&apartment.HasLaundry,
		&apartment.CreatedAt,
		&apartment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update apartment: %w", err)
	}

	return &apartment, nil
}

//go:embed delete.sql
var deleteApartmentQuery string

// DeleteApartment removes an apartment by ID
func (db *DB) DeleteApartment(id int64) error {

	result, err := db.Exec(deleteApartmentQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete apartment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("apartment with id %d not found", id)
	}

	return nil
}
