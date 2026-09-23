// Package model accept sqlite database
package model

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type schema struct {
	name    string
	content string
}

const ERR_DATABASE_EXIST = "database exists"

const credentialTableSchema = `
CREATE TABLE IF NOT EXISTS Credential (
  userAlias varchar(64) PRIMARY KEY,
  name varchar(64),
  cname varchar(64),
  password BLOB NOT NULL,
  role varchar(64) NOT NULL
);`

const studentTableSchema = `
CREATE TABLE IF NOT EXISTS Student (
  userAlias varchar(64) PRIMARY KEY,
	classCode TEXT,
	classNo INTEGER,
	priorities BLOB,
	isX3 BOOLEAN,
	isConfirmed BOOLEAN,
	ranking INTEGER DEFAULT 0,
	timestamp DATETIME NULL,
  FOREIGN KEY (userAlias) REFERENCES Credential(userAlias)
	);`

const subjectTableSchema = `
	CREATE TABLE IF NOT EXISTS Subject (
		code varchar(64) PRIMARY KEY,
		capacity INTEGER DEFAULT 0
	);`

const signatureTableSchema = `
	CREATE TABLE IF NOT EXISTS Signature (
		userAlias varchar(64) PRIMARY KEY,
  	isSigned BOOLEAN,
  	address TEXT,
		FOREIGN KEY (userAlias) REFERENCES Credential(userAlias)
	);`

var schemas = []schema{
	{"CREDENTIAL", credentialTableSchema},
	{"STUDENT", studentTableSchema},
	{"SUBJECT", subjectTableSchema},
	{"SIGNATURE", signatureTableSchema},
}

// CreateDatabase create new database and write schema on it.
// To create long-lived db for server, use `sql.Open()` to open a database.
func CreateDatabase(dsn string, isOverWrite bool) error {
	// Check if the SQLite database file already exists
	if _, err := os.Stat(dsn); err == nil {
		if !isOverWrite {
			return errors.New(ERR_DATABASE_EXIST)
		}
		// If overwrite is true, remove the existing database file
		if err := os.Remove(dsn); err != nil {
			return err
		}
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Enable WAL mode and set busy timeout for reliability
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		return err
	}

	// Enable foreign key support
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	if err := createTable(db); err != nil {
		return err
	}

	return nil
}

// createTable create table for the schema of ScheduleSchema
func createTable(db *sql.DB) error {
	for _, schema := range schemas {
		if _, err := db.Exec(schema.content); err != nil {
			return err
		}
		fmt.Printf("Add %s schema.\n", schema.name)
	}
	return nil
}
