package dbmanager

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"ssgo-server/model"
	"ssgo-server/model/master"
)

type Manager struct {
	mu           sync.RWMutex
	masterDB     *master.DB
	cohortDBs    map[string]*sql.DB
	dataDir      string
	activeCohort string
}

func New(dataDir string) (*Manager, error) {
	masterDSN := filepath.Join(dataDir, "master.sqlite")
	masterDB, err := master.OpenMasterDB(masterDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open master database: %v", err)
	}

	return &Manager{
		masterDB:  masterDB,
		cohortDBs: make(map[string]*sql.DB),
		dataDir:   dataDir,
	}, nil
}

func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	_ = m.masterDB.Close()
	for _, db := range m.cohortDBs {
		_ = db.Close()
	}
}

// Master returns the master database handle.
func (m *Manager) Master() *master.DB {
	return m.masterDB
}

// GetActiveDB returns the database connection for the currently active cohort.
func (m *Manager) GetActiveDB() (*sql.DB, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	active, err := m.masterDB.GetActiveCohort()
	if err != nil {
		return nil, "", err
	}

	m.activeCohort = active.ID

	// Check if already opened
	if db, exists := m.cohortDBs[active.ID]; exists {
		return db, active.ID, nil
	}

	// Open cohort DB
	dbPath := filepath.Join(m.dataDir, active.DBFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open cohort database %s: %v", active.ID, err)
	}

	// Enable WAL mode and set busy timeout
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, "", err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		db.Close()
		return nil, "", err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, "", err
	}

	m.cohortDBs[active.ID] = db
	return db, active.ID, nil
}

// GetCohortDB returns a database connection for a specific cohort.
func (m *Manager) GetCohortDB(id string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the dbFileName from master
	cohorts, err := m.masterDB.ListCohorts()
	if err != nil {
		return nil, err
	}

	var dbFileName string
	for _, c := range cohorts {
		if c.ID == id {
			dbFileName = c.DBFileName
			break
		}
	}

	if dbFileName == "" {
		return nil, master.ErrCohortNotFound
	}

	// Check if already opened
	if db, exists := m.cohortDBs[id]; exists {
		return db, nil
	}

	// Open cohort DB
	dbPath := filepath.Join(m.dataDir, dbFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cohort database %s: %v", id, err)
	}

	// Enable WAL mode and set busy timeout
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		db.Close()
		return nil, err
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, err
	}

	m.cohortDBs[id] = db
	return db, nil
}

// CreateCohortDB initializes a new physical SQLite database for a cohort and writes schemas.
func (m *Manager) CreateCohortDB(c *master.Cohort) error {
	dbPath := filepath.Join(m.dataDir, c.DBFileName)

	// Prevent overwriting physical file if it exists
	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("cohort database file %s already exists on disk", c.DBFileName)
	}

	// Ensure folder exists
	if err := os.MkdirAll(m.dataDir, 0755); err != nil {
		return err
	}

	// Use CreateDatabase function to handle the physical file and schemas
	if err := model.CreateDatabase(dbPath, true); err != nil {
		return err
	}

	// Create master list record
	return m.masterDB.CreateCohort(c)
}

// DeleteCohortDB permanently deletes a cohort record and its physical .sqlite file.
func (m *Manager) DeleteCohortDB(id string) error {
	m.mu.Lock()

	// Close cached connection if opened
	if db, exists := m.cohortDBs[id]; exists {
		_ = db.Close()
		delete(m.cohortDBs, id)
	}
	m.mu.Unlock()

	// Delete from master DB
	dbFileName, err := m.masterDB.DeleteCohort(id)
	if err != nil {
		return err
	}

	// Delete file from disk
	dbPath := filepath.Join(m.dataDir, dbFileName)
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Remove(dbPath); err != nil {
			return fmt.Errorf("deleted cohort metadata but failed to delete file %s from disk: %v", dbFileName, err)
		}
	}

	return nil
}
