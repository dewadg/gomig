package gomig

import (
	"strconv"
	"strings"
)

// Table represents slices of ColumnInterface types.
type Table struct {
	Columns []ColumnInterface
}

// ColumnInterface sets mandatory methods for Column-type.
type ColumnInterface interface {
	ToString() string
	IsPK() bool
}

// NewTable returns new Columns instance.
func NewTable() *Table {
	return &Table{Columns: make([]ColumnInterface, 0)}
}

// ToString returns all columns into SQL query string.
func (cols *Table) ToString() string {
	columns := make([]string, 0)
	for _, c := range cols.Columns {
		columns = append(columns, c.ToString())
	}

	return strings.Join(columns[:], ", ")
}

// IntegerCol represents an integer-type column.
type IntegerCol struct {
	Name       string
	IsPrimary  bool
	IsUnsigned bool
	IsNotNull  bool
}

// ToString returns current IntegerCol instance as a string
func (c *IntegerCol) ToString() string {
	output := c.Name + " int"

	if c.IsUnsigned {
		output += " UNSIGNED"
	}
	if c.IsPrimary {
		output += " NOT NULL PRIMARY KEY AUTO_INCREMENT"
	} else {
		if c.IsNotNull {
			output += " NOT NULL"
		}
	}

	return output
}

// IsPK determines if current IntegerCol is a primary key.
func (c *IntegerCol) IsPK() bool {
	return c.IsPrimary
}

// Primary sets current IntegerCol instance as primary key.
func (c *IntegerCol) Primary() *IntegerCol {
	c.IsPrimary = true
	return c
}

// NotNull sets current IntegerCol instance as not-null field.
func (c *IntegerCol) NotNull() *IntegerCol {
	c.IsNotNull = true
	return c
}

// Unsigned sets current IntegerCol instance as an unsigned field.
func (c *IntegerCol) Unsigned() *IntegerCol {
	c.IsUnsigned = true
	return c
}

// Integer returns new instance of IntegerCol.
func (cols *Table) Integer(name string) *IntegerCol {
	c := &IntegerCol{Name: name}
	cols.Columns = append(cols.Columns, c)
	return c
}

// VarcharCol represents an varchar-type column.
type VarcharCol struct {
	Name      string
	Length    int
	IsPrimary bool
	IsNotNull bool
}

// ToString returns current IntegerCol instance as a string
func (c *VarcharCol) ToString() string {
	output := c.Name

	if c.Length != 0 {
		output += " varchar(" + strconv.Itoa(c.Length) + ")"
	} else {
		output += " varchar"
	}
	if c.IsPrimary {
		output += " NOT NULL PRIMARY KEY"
	} else {
		if c.IsNotNull {
			output += " NOT NULL"
		}
	}

	return output
}

// IsPK determines if current VarcharCol is a primary key.
func (c *VarcharCol) IsPK() bool {
	return c.IsPrimary
}

// Primary sets current VarcharCol instance as primary key.
func (c *VarcharCol) Primary() *VarcharCol {
	c.IsPrimary = true
	return c
}

// NotNull sets current VarcharCol instance as not-null field.
func (c *VarcharCol) NotNull() *VarcharCol {
	c.IsNotNull = true
	return c
}

// Varchar returns new instance of VarcharCol.
func (cols *Table) Varchar(name string, length int) *VarcharCol {
	c := &VarcharCol{Name: name, Length: length}
	cols.Columns = append(cols.Columns, c)
	return c
}

// DateTimeCol represents an datetime-type column.
type DateTimeCol struct {
	Name      string
	IsNotNull bool
}

// ToString returns current DateTimeCol instance as a string
func (c *DateTimeCol) ToString() string {
	output := c.Name + " DATETIME"

	return output
}

// IsPK determines if current VarcharCol is a primary key.
func (c *DateTimeCol) IsPK() bool {
	return false
}

// NotNull sets current VarcharCol instance as not-null field.
func (c *DateTimeCol) NotNull() *DateTimeCol {
	c.IsNotNull = true
	return c
}

// DateTime returns new instance of DateTimeCol.
func (cols *Table) DateTime(name string) *DateTimeCol {
	c := &DateTimeCol{Name: name}
	cols.Columns = append(cols.Columns, c)
	return c
}
