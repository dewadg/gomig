package gomig

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	// Imports MySQL driver.
	_ "github.com/go-sql-driver/mysql"
)

func TestMigration(t *testing.T) {
	connString := "root:@tcp(localhost:3306)/dev_migrator"
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}

	m := New()
	m.SetDriver(db)
	m.SetMigrations(func() []*Migration {
		return []*Migration{
			&Migration{
				Name: "create_roles_table",
				Up: func() error {
					return m.Create("roles", func() *Table {
						table := NewTable()
						table.Integer("id").NotNull().Primary()
						table.Varchar("name", 255).NotNull()

						return table
					})
				},
				Down: func() error {
					return m.Drop("roles")
				},
			},
			&Migration{
				Name: "create_users_table",
				Up: func() error {
					return m.Create("users", func() *Table {
						table := NewTable()
						table.Integer("id").NotNull().Primary()
						table.Varchar("email", 255).NotNull()
						table.Varchar("password", 60).NotNull()
						table.DateTime("created_at").NotNull()
						table.Integer("role_id").NotNull()

						return table
					})
				},
				Down: func() error {
					return m.Drop("users")
				},
			},
			&Migration{
				Name: "create_posts_table",
				Up: func() error {
					return m.Create("posts", func() *Table {
						table := NewTable()
						table.Integer("id").NotNull().Primary()
						table.Varchar("title", 255).NotNull()
						table.Enum("type", "post", "page", "static").NotNull()

						return table
					})
				},
				Down: func() error {
					return m.Drop("posts")
				},
			},
		}
	})

	n, err := m.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d migration(s) succeeded\n", n)
}
