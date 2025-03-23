package gomb

import (
	"fmt"
	"strings"
)

// Table represents a database table
type Table struct {
	Name       string         `json:"name"`
	Label      string         `json:"label"`
	Columns    []*Column      `json:"columns"`
	Attributes map[string]any `json:"attributes"`
	Comment    string         `json:"comment"`
}

// NewTable initializes and returns a new Table instance
func NewTable(name string) *Table {
	return &Table{Name: name}
}

// AddColumn adds a new column to the table
func (t *Table) AddColumn(column *Column) *Table {
	t.Columns = append(t.Columns, column)
	return t
}

// TableDefinition generates the full SQL table definition by combining various table attributes
func (t *Table) ToSQL() (string, []error) {
	var def []string
	var errors []error

	// Add table name
	if t.Name == "" {
		errors = append(errors, fmt.Errorf("table name cannot be empty"))
		return "", errors
	}
	def = append(def, fmt.Sprintf("CREATE TABLE %s", t.Name))

	// Add columns
	columnDefs := []string{}
	for _, col := range t.Columns {
		colSQL, err := col.ToSQL()
		if err != nil {
			errors = append(errors, err)
			continue // Skip this column if there's an error
		}
		columnDefs = append(columnDefs, colSQL)
	}

	if len(columnDefs) == 0 {
		errors = append(errors, fmt.Errorf("no valid columns defined for table %s", t.Name))
		return "", errors
	}

	def = append(def, fmt.Sprintf("(%s)", strings.Join(columnDefs, ", ")))

	// Add table-level comment if provided
	if t.Comment != "" {
		def = append(def, fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", t.Name, t.Comment))
	}

	if len(errors) > 0 {
		return "", errors
	}

	// Join and return the SQL definition
	return strings.Join(def, " "), nil
}

// Validate validates the table and its columns
func (t *Table) Validate() []error {
	var errors []error

	// Check if table name is empty
	if t.Name == "" {
		errors = append(errors, fmt.Errorf("table name cannot be empty"))
	}

	return errors
}
