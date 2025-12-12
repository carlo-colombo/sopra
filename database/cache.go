package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/carlo-colombo/sopra/migrations"
	"github.com/carlo-colombo/sopra/model"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

// DB handles the database operations for caching.
type DB struct {
	db *sql.DB
}

// runMigrations applies the database migrations.
func runMigrations(dataSourceName string) error {
	d, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, "sqlite3://"+dataSourceName)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// NewDB initializes the SQLite database and returns a DB instance.
func NewDB(dataSourceName string) (*DB, error) {
	if err := runMigrations(dataSourceName); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Close closes the underlying database connection.
func (c *DB) Close() error {
	return c.db.Close()
}

// GetFlightCount returns the total number of flights in the cache.
func (c *DB) GetFlightCount() (int, error) {
	var count int
	err := c.db.QueryRow("SELECT COUNT(*) FROM flight_log").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetFlight retrieves a cached FlightInfo by key.
func (c *DB) GetFlight(key string) (*model.FlightInfo, time.Time, error) {
	var jsonValue string
	var lastSeen time.Time
	var identificationCount int
	err := c.db.QueryRow("SELECT value, last_seen, identification_count FROM flight_log WHERE key = ?", key).Scan(&jsonValue, &lastSeen, &identificationCount)
	if err == sql.ErrNoRows {
		return nil, time.Time{}, nil // Cache miss
	}
	if err != nil {
		return nil, time.Time{}, err
	}

	var flightInfo model.FlightInfo
	if err := json.Unmarshal([]byte(jsonValue), &flightInfo); err != nil {
		return nil, time.Time{}, err
	}
	flightInfo.IdentificationCount = identificationCount

	log.Printf("Cache hit for key: %s, last seen: %s\n", key, lastSeen)
	return &flightInfo, lastSeen, nil
}

// GetAllFlights retrieves all the logged FlightInfo.
func (c *DB) GetAllFlights() ([]*model.FlightInfo, []time.Time, error) {
	rows, err := c.db.Query("SELECT value, last_seen, identification_count FROM flight_log ORDER BY last_seen DESC")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var flights []*model.FlightInfo
	var lastSeens []time.Time

	for rows.Next() {
		var jsonValue string
		var lastSeen time.Time
		var identificationCount int
		if err := rows.Scan(&jsonValue, &lastSeen, &identificationCount); err != nil {
			return nil, nil, err
		}

		var flightInfo model.FlightInfo
		if err := json.Unmarshal([]byte(jsonValue), &flightInfo); err != nil {
			return nil, nil, err
		}
		flightInfo.IdentificationCount = identificationCount
		flights = append(flights, &flightInfo)
		lastSeens = append(lastSeens, lastSeen)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return flights, lastSeens, nil
}

// LogFlight stores a FlightInfo in the cache.
func (c *DB) LogFlight(key string, flightInfo *model.FlightInfo) error {
	jsonValue, err := json.Marshal(flightInfo)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("INSERT INTO flight_log (key, value, last_seen, identification_count) VALUES (?, ?, ?, 1) ON CONFLICT(key) DO UPDATE SET value = excluded.value, last_seen = excluded.last_seen, identification_count = identification_count + 1", key, string(jsonValue), time.Now())
	log.Printf("Logged flight for key: %s\n", key)
	return err
}

// GetLatestFlight retrieves the most recently logged FlightInfo.
func (c *DB) GetLatestFlight() (*model.FlightInfo, time.Time, error) {
	var jsonValue string
	var lastSeen time.Time
	var identificationCount int
	err := c.db.QueryRow("SELECT value, last_seen, identification_count FROM flight_log ORDER BY last_seen DESC LIMIT 1").Scan(&jsonValue, &lastSeen, &identificationCount)
	if err == sql.ErrNoRows {
		return nil, time.Time{}, nil // No flights in cache
	}
	if err != nil {
		return nil, time.Time{}, err
	}

	var flightInfo model.FlightInfo
	if err := json.Unmarshal([]byte(jsonValue), &flightInfo); err != nil {
		return nil, time.Time{}, err
	}
	flightInfo.IdentificationCount = identificationCount

	return &flightInfo, lastSeen, nil
}

// ClearFlightLog deletes all records from the flight_log table.
func (c *DB) ClearFlightLog() error {
	_, err := c.db.Exec("DELETE FROM flight_log")
	return err
}
