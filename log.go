package gomig

import "time"
import "database/sql"

// MigrationLog represents logs from migrations table.
type MigrationLog struct {
	ID    int64
	Name  string
	RunAt time.Time
}

// createMigrationsTable creates migrations table if not exists.
func (m *Migrator) createMigrationsTable() error {
	return m.Create("migrations", func() *Table {
		table := NewTable()
		table.Integer("id").Unsigned().Primary()
		table.Varchar("name", 255).NotNull()
		table.DateTime("runAt").NotNull()

		return table
	})
}

// fetchLogs fetches logs from migrations table and returns them as *MigrationLog instance.
func (m *Migrator) fetchLogs() error {
	var (
		fetchedID    sql.NullInt64
		fetchedName  sql.NullString
		fetchedRunAt sql.NullString
		logs         = make([]*MigrationLog, 0)
	)

	rows, err := m.driver.Query("SELECT id, name, runAt FROM migrations")
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(&fetchedID, &fetchedName, &fetchedRunAt); err != nil {
			return err
		}

		runAt, err := time.Parse("2006-01-02 15:04:05", fetchedRunAt.String)
		if err != nil {
			return err
		}
		logs = append(logs, &MigrationLog{ID: fetchedID.Int64, Name: fetchedName.String, RunAt: runAt})
	}

	m.logs = logs
	return nil
}

// log logs current migration status to migrations table.
func (m *Migrator) log(mgrn *Migration) error {
	currentTime := time.Now()
	stmt, err := m.driver.Prepare("INSERT INTO migrations SET name = ?, runAt = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(mgrn.Name, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.logs = append(m.logs, &MigrationLog{ID: logID, Name: mgrn.Name, RunAt: currentTime})
	return err
}

// hasLog determines if the given migration already run before.
func (m *Migrator) hasLog(migration string) bool {
	for _, l := range m.logs {
		if l.Name == migration {
			return true
		}
	}
	return false
}
