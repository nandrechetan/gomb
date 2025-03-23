package gomb

import (
	"fmt"
)

type DropTable struct {
	Name    string
	Cascade bool // If true, adds CASCADE to the DROP statement
}

// NewTable initializes and returns a new Table instance
func NewDropTable(name string) *DropTable {
	return &DropTable{Name: name}
}

// SetCascade enables or disables the CASCADE option
func (t *DropTable) SetCascade(cascade bool) *DropTable {
	t.Cascade = cascade
	return t
}

// ToSQL generates the DROP TABLE SQL statement
func (t *DropTable) ToSQL() (string, error) {
	if t.Name == "" {
		return "", fmt.Errorf("table name cannot be empty")
	}

	// Construct DROP TABLE statement
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", t.Name)
	if t.Cascade {
		sql += " CASCADE"
	}

	return sql, nil
}
