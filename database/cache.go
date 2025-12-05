package database

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/carlo-colombo/sopra/model"
)

// Cache handles the database operations for caching.
type Cache struct {
	db *sql.DB
}

// NewCache initializes the SQLite database and returns a Cache instance.
func NewCache(dataSourceName string) (*Cache, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the cache table if it doesn't exist.
	// The value is stored as TEXT and will contain the JSON response.
	_ , err = db.Exec(`CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value TEXT
	)`)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

// Get retrieves a cached FlightInfo by key.
func (c *Cache) Get(key string) (*model.FlightInfo, error) {
	var jsonValue string
	err := c.db.QueryRow("SELECT value FROM cache WHERE key = ?", key).Scan(&jsonValue)
	if err == sql.ErrNoRows {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var flightInfo model.FlightInfo
	if err := json.Unmarshal([]byte(jsonValue), &flightInfo); err != nil {
		return nil, err
	}

	log.Printf("Cache hit for key: %s\n", key)
	return &flightInfo, nil
}

// Set stores a FlightInfo in the cache.
func (c *Cache) Set(key string, flightInfo *model.FlightInfo) error {
	jsonValue, err := json.Marshal(flightInfo)
	if err != nil {
		return err
	}

	_ , err = c.db.Exec("INSERT OR REPLACE INTO cache (key, value) VALUES (?, ?)", key, string(jsonValue))
	log.Printf("Cached value for key: %s\n", key)
	return err
}
