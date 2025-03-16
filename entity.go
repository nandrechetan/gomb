package gomb

// Table represents a database table
type Table struct {
	Name       string         `json:"name"`
	Label      string         `json:"label"`
	Columns    []*Column      `json:"columns"`
	Attributes map[string]any `json:"attributes"`
	Comment    string         `json:"comment"`
}

// AlterTable represents an ALTER TABLE statement
type AlterTable struct {
	TableName  string
	Operations []ColumnOperation
	Comment    string
}

// Column represents a database column
type Column struct {
	Name             string         `json:"name"`        // Column name
	DataType         DataType       `json:"data_type"`   // Data type (e.g., VARCHAR, INTEGER, etc.)
	Length           int            `json:"length"`      // Length (e.g., VARCHAR(255)
	Precision        int            `json:"precision"`   // For DECIMAL or numeric types
	Scale            int            `json:"scale"`       // For DECIMAL or numeric types
	PrimaryKey       bool           `json:"primary_key"` // Whether this column is a primary key
	AutoNumber       bool           `json:"auto_number"` // Whether this column is auto-incrementing
	AutoNumberStart  int            `json:"auto_number_start"`
	AutoNumberPrefix string         `json:"auto_number_prefix"`
	NotNull          bool           `json:"not_null"`       // Whether this column allows NULL values
	Unique           bool           `json:"unique"`         // Whether this column has a UNIQUE constraint
	Default          string         `json:"default"`        // Default value for the column
	Check            string         `json:"check"`          // CHECK constraint expression
	References       string         `json:"references"`     // Foreign key reference (e.g., "other_table(column)")
	Generated        string         `json:"generated"`      // Expression for generated columns
	Collation        string         `json:"collation"`      // Collation for text-based columns
	Comment          string         `json:"comment"`        // Comment or description of the column
	Storage          string         `json:"storage"`        // Storage option (e.g., PLAIN, EXTERNAL)
	Compression      string         `json:"compression"`    // Compression method (PostgreSQL 14+)
	IdentityStart    int            `json:"identity_start"` // Start value for identity columns
	IdentityInc      int            `json:"identity_inc"`   // Increment value for identity columns
	Attributes       map[string]any `json:"attributes"`     // Custom/extensible attributes
	NewName          string         `json:"new_name"`
	NewDataType      DataType       `json:"new_data_type"`
}

// ColumnOperation represents a single operation on a column
type ColumnOperation struct {
	Operation AlterTableOperation
	Column    *Column
}

// Index represents a database index
type Index struct {
	name           string
	table          string
	columns        []string
	unique         bool
	concurrently   bool
	using          string
	where          string
	schema         string
	includeColumns []string
	method         string // btree, hash, gist, gin, etc.
	tablespace     string
	withOptions    []string
}
