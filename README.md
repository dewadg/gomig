## gomig

gomig is a Go library which helps you handle database migrations. Heavily inspired by Laravel's migrations.

## Supported Database
Currently supports only: MariaDB, MySQL

## Installation
```go
go get github.com/dewadg/gomig
```

## Usage
Import gomig to your project:
```go
import "github.com/dewadg/gomig"
```

Create a gomig instance:
```go
migrator := gomig.New()
```

Sets a driver:
```go
db, err := sql.Open("mysql", "root:boom@tcp(localhost:3306)/dev_migrator")
if err != nil {
    log.Fatal(err)
}

migrator.SetDriver(db)
```

Migrations are defined as an array of *gomig.Migration instances:
```go
migrator.SetMigrations(func() []*gomig.Migration {
    return []*gomig.Migration{
        &gomig.Migration{
            Name: "create_roles_table",
            Up: func() error {
                return migrator.Create("roles", func() *gomig.Table {
                    table := gomig.NewTable()
                    table.Integer("id").NotNull().Primary()
                    table.Varchar("name", 255).NotNull()

                    return table
                })
            },
            Down: func() error {
                return migrator.Drop("roles")
            },
        },
        &gomig.Migration{
            Name: "create_users_table",
            Up: func() error {
                return migrator.Create("users", func() *gomig.Table {
                    table := gomig.NewTable()
                    table.Integer("id").NotNull().Primary()
                    table.Varchar("email", 255).NotNull()
                    table.Varchar("password", 60).NotNull()
                    table.DateTime("created_at").NotNull()
                    table.Integer("role_id").NotNull()

                    return table
                })
            },
            Down: func() error {
                return migrator.Drop("users")
            },
        },
    }
})
```

Run the migrations:
```go
n, err := migrator.Migrate()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%d migration(s) succeeded\n", n)
```

## Available Column-defining Methods

#### Integer
```go
table.Integer("name")
```

#### Varchar
```go
table.Varchar("name", 255) // 255 is the length of the column
```

#### DateTime
```go
table.DateTime("name")
```

### Other Methods
Column defining methods can also have other methods to be chained for several use:
```go
table.Integer("name").Unsigned().NotNull()
```

```go
table.Integer("name").NotNull()
```

```go
table.Integer("name").Unsigned().Primary() // Only 1 column should be defined as Primary Key
```

## To-do

1. Rollback capability
2. Adding column-modifying methods
3. Relationship constraints 
4. Adding more column-defining methods
