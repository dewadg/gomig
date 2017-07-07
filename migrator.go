package gomig

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// ErrDuplicatePrimaryKey represents an error message where more than 1 primary key were set
	ErrDuplicatePrimaryKey = "More than 1 primary key were set"
)

// Migration represents a migration.
type Migration struct {
	Name string
	Up   func() error
	Down func() error
}

// Migrator represents a migrator tool instance.
type Migrator struct {
	migrations []*Migration
	driver     *sql.DB
	logs       []*MigrationLog
}

// New returns new instance of Migrator.
func New() *Migrator {
	return &Migrator{logs: make([]*MigrationLog, 0)}
}

// SetDriver sets the database driver to run migrations on.
func (m *Migrator) SetDriver(driver *sql.DB) {
	m.driver = driver
}

// SetMigrations sets migrations to the current Migrator.
func (m *Migrator) SetMigrations(f func() []*Migration) {
	m.migrations = f()
}

// Migrate runs migration on the current Migrator.
func (m *Migrator) Migrate() (int, error) {
	var successMigrations int

	m.createMigrationsTable()
	if err := m.fetchLogs(); err != nil {
		return successMigrations, err
	}
	for _, mgrn := range m.migrations {
		if m.hasLog(mgrn.Name) {
			continue
		}

		if err := mgrn.Up(); err != nil {
			return successMigrations, err
		}

		m.log(mgrn)
		successMigrations++
		fmt.Println(time.Now().Format("2006-01-02 15:04:09") + " " + mgrn.Name)
	}
	return successMigrations, nil
}

// Create creates new table based on defined columns.
func (m *Migrator) Create(name string, f func() *Table) error {
	var foundPK int
	var columns []string

	// Ensure only 1 column set as primary key
	// and converts *Table to array of strings.
	for _, c := range f().Columns {
		if c.IsPK() {
			foundPK++
			if foundPK > 1 {
				return errors.New(ErrDuplicatePrimaryKey)
			}
		}
		columns = append(columns, c.ToString())
	}

	// Convert data to SQL query string and execute it.
	queryString := []string{"CREATE TABLE ", name, " (", strings.Join(columns[:], ", "), ")"}
	_, err := m.driver.Exec(strings.Join(queryString, ""))

	if err != nil {
		return err
	}
	return nil
}

// Drop drops given table.
func (m *Migrator) Drop(name string) error {
	sql := "DROP TABLE " + name
	_, err := m.driver.Exec(sql)

	if err != nil {
		return err
	}
	return nil
}
