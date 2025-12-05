package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/carlo-colombo/sopra/model"
)

// DB handles the database operations for caching.
type DB struct {
	db *sql.DB
}

// NewDB initializes the SQLite database and returns a DB instance.
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the flight_log table if it doesn't exist.
	// The value is stored as TEXT and will contain the JSON response.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS flight_log (
		key TEXT PRIMARY KEY,
		value TEXT,
		last_seen DATETIME
	)`)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Close closes the underlying database connection.
func (c *DB) Close() error {
	return c.db.Close()
}



// GetFlight retrieves a cached FlightInfo by key.
func (c *DB) GetFlight(key string) (*model.FlightInfo, time.Time, error) {
	var jsonValue string
	var lastSeen time.Time
	err := c.db.QueryRow("SELECT value, last_seen FROM flight_log WHERE key = ?", key).Scan(&jsonValue, &lastSeen)
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

	log.Printf("Cache hit for key: %s, last seen: %s\n", key, lastSeen)
	return &flightInfo, lastSeen, nil
}

// LogFlight stores a FlightInfo in the cache.
func (c *DB) LogFlight(key string, flightInfo *model.FlightInfo) error {
	jsonValue, err := json.Marshal(flightInfo)
	if err != nil {
		return err
	}

	_, err = c.db.Exec("INSERT OR REPLACE INTO flight_log (key, value, last_seen) VALUES (?, ?, ?)", key, string(jsonValue), time.Now())
	log.Printf("Logged flight for key: %s\n", key)
	return err
}

// GetLatestFlight retrieves the most recently logged FlightInfo.
func (c *DB) GetLatestFlight() (*model.FlightInfo, error) {
	var jsonValue string
	err := c.db.QueryRow("SELECT value FROM flight_log ORDER BY last_seen DESC LIMIT 1").Scan(&jsonValue)
	if err == sql.ErrNoRows {
		return nil, nil // No flights in cache
	}
	if err != nil {
		return nil, err
	}

	var flightInfo model.FlightInfo
	if err := json.Unmarshal([]byte(jsonValue), &flightInfo); err != nil {
		return nil, err
	}

	return &flightInfo, nil
}

