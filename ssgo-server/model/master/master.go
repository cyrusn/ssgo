package master

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCohortNotFound    = errors.New("cohort not found")
	ErrNoActiveCohort    = errors.New("no active cohort configured")
	ErrSuperadminNotFound = errors.New("superadmin not found")
	ErrInvalidPassword   = errors.New("invalid password")
)

type DB struct {
	*sql.DB
}

type Superadmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Cohort struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DBFileName string `json:"dbFileName"`
	IsActive   bool   `json:"isActive"`
	Config     string `json:"config"` // JSON string representation of frontend configs
}

const superadminTableSchema = `
CREATE TABLE IF NOT EXISTS Superadmin (
	username TEXT PRIMARY KEY,
	password BLOB NOT NULL
);`

const cohortTableSchema = `
CREATE TABLE IF NOT EXISTS Cohort (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	dbFileName TEXT NOT NULL,
	isActive BOOLEAN DEFAULT 0,
	config TEXT NOT NULL
);`

// OpenMasterDB opens or creates the master.sqlite database.
func OpenMasterDB(masterDSN string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(masterDSN)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", masterDSN)
	if err != nil {
		return nil, err
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

	if _, err := db.Exec(superadminTableSchema); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(cohortTableSchema); err != nil {
		db.Close()
		return nil, err
	}

	masterDB := &DB{db}
	if err := masterDB.SeedSuperadmins(); err != nil {
		db.Close()
		return nil, err
	}

	return masterDB, nil
}

// SeedSuperadmins inserts default superadmins if they do not exist.
func (db *DB) SeedSuperadmins() error {
	username := os.Getenv("SUPER_ADMIN_USERNAME")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if username == "" || password == "" {
		fmt.Println("Warning: SUPER_ADMIN_USERNAME or SUPER_ADMIN_PASSWORD not set in environment.")
		return nil
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Superadmin WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO Superadmin (username, password) VALUES (?, ?)", username, hashed)
		if err != nil {
			return err
		}
		fmt.Printf("Seeded root Superadmin user: %s\n", username)
	}
	return nil
}

// ListSuperadmins returns all superadmin usernames.
func (db *DB) ListSuperadmins() ([]string, error) {
	rows, err := db.Query("SELECT username FROM Superadmin")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		admins = append(admins, u)
	}
	return admins, nil
}

// CreateSuperadmin adds a new superadmin account.
func (db *DB) CreateSuperadmin(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO Superadmin (username, password) VALUES (?, ?)", username, hashed)
	return err
}

// UpdateSuperadminPassword changes a superadmin's password.
func (db *DB) UpdateSuperadminPassword(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := db.Exec("UPDATE Superadmin SET password = ? WHERE username = ?", hashed, username)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSuperadminNotFound
	}
	return nil
}

// DeleteSuperadmin removes a superadmin account, protecting the root from .env.
func (db *DB) DeleteSuperadmin(username string) error {
	rootUsername := os.Getenv("SUPER_ADMIN_USERNAME")
	if username == rootUsername {
		return errors.New("cannot delete the root superadmin account")
	}

	res, err := db.Exec("DELETE FROM Superadmin WHERE username = ?", username)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSuperadminNotFound
	}
	return nil
}

// ResetRootPasswordToDefault resets the root user's password to the value in .env.
func (db *DB) ResetRootPasswordToDefault() error {
	username := os.Getenv("SUPER_ADMIN_USERNAME")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if username == "" || password == "" {
		return errors.New("SUPER_ADMIN_USERNAME or SUPER_ADMIN_PASSWORD not set in environment")
	}

	return db.UpdateSuperadminPassword(username, password)
}

// Authenticate checks superadmin credentials.
func (db *DB) Authenticate(username, password string) error {
	var hashedPassword []byte
	err := db.QueryRow("SELECT password FROM Superadmin WHERE username = ?", username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return ErrSuperadminNotFound
	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// CreateCohort creates a new cohort record and initializes its database.
func (db *DB) CreateCohort(c *Cohort) error {
	// Check if cohort already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Cohort WHERE id = ?", c.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cohort with ID '%s' already exists", c.ID)
	}

	// Insert into master
	_, err = db.Exec(
		"INSERT INTO Cohort (id, name, dbFileName, isActive, config) VALUES (?, ?, ?, ?, ?)",
		c.ID, c.Name, c.DBFileName, false, c.Config,
	)
	return err
}

// ListCohorts lists all configured cohorts.
func (db *DB) ListCohorts() ([]Cohort, error) {
	rows, err := db.Query("SELECT id, name, dbFileName, isActive, config FROM Cohort")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cohorts := []Cohort{}
	for rows.Next() {
		var c Cohort
		err := rows.Scan(&c.ID, &c.Name, &c.DBFileName, &c.IsActive, &c.Config)
		if err != nil {
			return nil, err
		}
		cohorts = append(cohorts, c)
	}
	return cohorts, nil
}

// SetCohortActive activates a specific cohort and deactivates all others.
func (db *DB) SetCohortActive(id string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deactivate all
	_, err = tx.Exec("UPDATE Cohort SET isActive = 0")
	if err != nil {
		return err
	}

	// Activate specified
	res, err := tx.Exec("UPDATE Cohort SET isActive = 1 WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCohortNotFound
	}

	return tx.Commit()
}

// DeactivateCohort deactivates a specific cohort.
func (db *DB) DeactivateCohort(id string) error {
	res, err := db.Exec("UPDATE Cohort SET isActive = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCohortNotFound
	}
	return nil
}

// GetActiveCohort returns the currently active cohort.
func (db *DB) GetActiveCohort() (*Cohort, error) {
	var c Cohort
	err := db.QueryRow("SELECT id, name, dbFileName, isActive, config FROM Cohort WHERE isActive = 1").
		Scan(&c.ID, &c.Name, &c.DBFileName, &c.IsActive, &c.Config)
	if err == sql.ErrNoRows {
		return nil, ErrNoActiveCohort
	} else if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpdateCohortConfig updates a cohort's frontend configuration JSON.
func (db *DB) UpdateCohortConfig(id string, configJSON string) error {
	// Validate JSON
	var temp map[string]interface{}
	if err := json.Unmarshal([]byte(configJSON), &temp); err != nil {
		return fmt.Errorf("invalid config JSON: %v", err)
	}

	res, err := db.Exec("UPDATE Cohort SET config = ? WHERE id = ?", configJSON, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCohortNotFound
	}
	return nil
}

// UpdateCohortName updates a cohort's display name.
func (db *DB) UpdateCohortName(id string, name string) error {
	res, err := db.Exec("UPDATE Cohort SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrCohortNotFound
	}
	return nil
}

// DeleteCohort deletes a cohort from the master list.
func (db *DB) DeleteCohort(id string) (string, error) {
	var dbFileName string
	err := db.QueryRow("SELECT dbFileName FROM Cohort WHERE id = ?", id).Scan(&dbFileName)
	if err == sql.ErrNoRows {
		return "", ErrCohortNotFound
	} else if err != nil {
		return "", err
	}

	_, err = db.Exec("DELETE FROM Cohort WHERE id = ?", id)
	if err != nil {
		return "", err
	}

	return dbFileName, nil
}
