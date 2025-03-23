package gomb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	// NewName          string         `json:"new_name"`
	// NewDataType      DataType       `json:"new_data_type"`
	UpdateOptions *ColumnUpdate `json:"update_options,omitempty"`
}

// ColumnUpdate holds modification details for a column
type ColumnUpdate struct {
	Name     string   `json:"name,omitempty"`
	DataType DataType `json:"data_type,omitempty"`
}

// NewTable initializes and returns a new Table instance
func NewColumn(name string) *Column {
	return &Column{Name: name}
}

// func (c *Column) SetNewName(name string) *Column {
// 	c.NewName = name
// 	return c
// }
// func (c *Column) SetNewDataType(name DataType) *Column {
// 	c.NewDataType = name
// 	return c
// }

// SetNewName sets a new name for the column
func (c *Column) SetNewName(newName string) *Column {
	if c.UpdateOptions == nil {
		c.UpdateOptions = &ColumnUpdate{}
	}
	c.UpdateOptions.Name = newName
	return c
}

// SetNewDataType sets a new data type for the column
func (c *Column) SetNewDataType(newDataType DataType) *Column {
	if c.UpdateOptions == nil {
		c.UpdateOptions = &ColumnUpdate{}
	}
	c.UpdateOptions.DataType = newDataType
	return c
}

// SetName sets the column name
func (c *Column) SetName(name string) *Column {
	c.Name = name
	return c
}

// SetPrimaryKey marks the column as a primary key
func (c *Column) SetPrimaryKey() *Column {
	c.PrimaryKey = true
	return c
}

// SetUnique marks the column as unique
func (c *Column) SetUnique() *Column {
	c.Unique = true
	return c
}

// SetNotNull marks the column as not allowing null values
func (c *Column) SetNotNull() *Column {
	c.NotNull = true
	return c
}

// SetDataType sets the data type of the column
func (c *Column) SetDataType(dataType DataType) *Column {
	c.DataType = dataType
	return c
}

// SetDefault sets the default value for the column
func (c *Column) SetDefault(defaultValue any) *Column {
	switch v := defaultValue.(type) {
	case string:
		c.Default = v
	case int, int32, int64, float32, float64:
		// For numeric values, convert to string without quotes
		c.Default = fmt.Sprintf("%v", v)
	case bool:
		// For boolean values, convert to lowercase string without quotes
		c.Default = fmt.Sprintf("%t", v)
	case nil:
		c.Default = "NULL"
	default:
		// For other types, convert to string and add quotes
		c.Default = fmt.Sprintf("%v", v)
	}
	return c
}

// AutoNumber generates the auto-number clause with a custom prefix
func (col *Column) SetAutoNumber() *Column {
	col.AutoNumber = true
	return col
}

// AutoNumberWithPrefixClause generates the auto-number clause with a custom prefix
func (col *Column) SetAutoNumberWithPrefix(startNumber int, prefix string) *Column {
	col.AutoNumber = true
	col.AutoNumberStart = startNumber
	col.AutoNumberPrefix = prefix
	return col
}
func T(tableName string) string {
	return tableName
}
func C(columnName string) string {
	return columnName
}

// SetReferences sets a foreign key reference
func (c *Column) SetReferences(table string, column string) *Column {
	c.References = fmt.Sprintf("%s(%s)", table, column)
	return c
}

// SetReferences sets a foreign key reference
func (c *Column) SetReferencesOnDeleteCascade(table string, column string) *Column {
	c.References = fmt.Sprintf("%s(%s) ON DELETE CASCADE", table, column)
	return c
}

func (col *Column) SetCheck(check string) *Column {
	col.Check = check
	return col
}
func (col *Column) SetGenerated(check string) *Column {
	col.Generated = check
	return col
}
func (col *Column) SetLength(length int) *Column {
	col.Length = length
	return col
}

func (col *Column) SetPrecision(precision int) *Column {
	col.Precision = precision
	return col
}

func (col *Column) SetScale(scale int) *Column {
	col.Scale = scale
	return col
}

func (col *Column) SetCollation(collation string) *Column {
	col.Collation = collation
	return col
}

func (col *Column) SetComment(comment string) *Column {
	col.Comment = comment
	return col
}

func (col *Column) SetStorage(storage string) *Column {
	col.Storage = storage
	return col
}

func (col *Column) SetCompression(compression string) *Column {
	col.Compression = compression
	return col
}

func (col *Column) SetIdentityStart(identityStart int) *Column {
	col.IdentityStart = identityStart
	return col
}

func (col *Column) SetIdentityIncrement(identityInc int) *Column {
	col.IdentityInc = identityInc
	return col
}

func (col *Column) SetAttributes(attributes map[string]any) *Column {
	col.Attributes = attributes
	return col
}

// ToSQL generates the SQL definition for the column
func (c *Column) ToSQL() (string, error) {
	// Validate the column definition first
	if err := c.Validate(); err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(c.Name)
	builder.WriteString(" ")

	// Add data type
	if c.DataType != "" {
		builder.WriteString(c.ToDataTypeString(c.DataType))
	}

	// Add primary key constraint
	if c.PrimaryKey {
		builder.WriteString(" PRIMARY KEY")
	}

	// Add auto-increment settings
	if c.AutoNumber {
		builder.WriteString(" AUTOINCREMENT")
		if c.AutoNumberStart > 0 {
			builder.WriteString(fmt.Sprintf(" START WITH %d", c.AutoNumberStart))
		}
		if c.AutoNumberPrefix != "" {
			builder.WriteString(fmt.Sprintf(" PREFIX '%s'", c.AutoNumberPrefix))
		}
	}

	// Add not null constraint
	if c.NotNull {
		builder.WriteString(" NOT NULL")
	}

	// Add unique constraint
	if c.Unique {
		builder.WriteString(" UNIQUE")
	}

	// Add default value
	if c.Default != "" {
		// Check if the default value needs quotes
		needsQuotes := true

		// Standard SQL functions/constants that don't need quotes
		switch c.Default {
		case "CURRENT_TIMESTAMP", "CURRENT_DATE", "CURRENT_TIME",
			"LOCAL_TIME", "LOCAL_TIMESTAMP", "TRUE", "FALSE", "NULL":
			needsQuotes = false
		}

		// Check if it's a numeric value
		if _, err := strconv.ParseFloat(c.Default, 64); err == nil {
			needsQuotes = false
		}

		// Boolean literals don't need quotes
		if c.Default == "true" || c.Default == "false" {
			needsQuotes = false
		}

		// Apply quotes if needed
		if needsQuotes {
			builder.WriteString(fmt.Sprintf(" DEFAULT '%s'", c.Default))
		} else {
			builder.WriteString(fmt.Sprintf(" DEFAULT %s", c.Default))
		}
	}

	// Add check constraint
	if c.Check != "" {
		builder.WriteString(fmt.Sprintf(" CHECK %s", c.Check))
	}

	// Add references (foreign key)
	if c.References != "" {
		builder.WriteString(fmt.Sprintf(" REFERENCES %s", c.References))
	}

	// Add generated column
	if c.Generated != "" {
		builder.WriteString(fmt.Sprintf(" GENERATED ALWAYS AS (%s)", c.Generated))
	}

	// Add collation
	if c.Collation != "" {
		builder.WriteString(fmt.Sprintf(" COLLATE %s", c.Collation))
	}

	// Add comment
	if c.Comment != "" {
		builder.WriteString(fmt.Sprintf(" COMMENT '%s'", c.Comment))
	}

	// Add storage option
	if c.Storage != "" {
		builder.WriteString(fmt.Sprintf(" STORAGE %s", c.Storage))
	}

	// Add compression method
	if c.Compression != "" {
		builder.WriteString(fmt.Sprintf(" COMPRESSION %s", c.Compression))
	}

	// Add identity settings
	if c.IdentityStart > 0 && c.IdentityInc > 0 {
		builder.WriteString(fmt.Sprintf(" IDENTITY (%d,%d)", c.IdentityStart, c.IdentityInc))
	}

	// Add custom attributes
	for key, value := range c.Attributes {
		builder.WriteString(fmt.Sprintf(" %s %v", key, value))
	}

	return builder.String(), nil
}

// Valid PostgreSQL data types for validation
var validDataTypes = map[DataType]bool{
	SerialType:   true,
	StringType:   true,
	IntegerType:  true,
	DecimalType:  true,
	BooleanType:  true,
	DateType:     true,
	DateTimeType: true,
}

func (col *Column) Validate() error {
	// Data Type Validation
	if !validDataTypes[col.DataType] {
		return fmt.Errorf("invalid data type: %s", col.DataType)
	}

	// Auto Number Start Validation
	if col.AutoNumber && col.AutoNumberStart < 0 {
		return errors.New("auto-number start must be greater or equal than 0")
	}

	// NotNull and Default Validation
	if col.NotNull && col.Default != "" {
		return errors.New("column cannot be both NOT NULL and have a DEFAULT value")
	}

	// IdentityStart and IdentityInc Validation
	if col.IdentityStart > 0 && col.IdentityInc <= 0 {
		return errors.New("identity increment must be greater than 0")
	}

	// Check constraint validation (if applicable)
	if col.Check != "" && !strings.Contains(col.Check, "(") {
		return errors.New("check constraint must have an expression in parentheses")
	}

	// References validation (foreign key format)
	if col.References != "" && !strings.Contains(col.References, "(") {
		return errors.New("foreign key references must be in the format 'table(column)'")
	}

	return nil
}

// Map the DataType to corresponding PostgreSQL data type
func (col *Column) ToDataType() string {
	return col.ToDataTypeString(col.DataType)
}
func (col *Column) ToNewDataType() string {
	return col.ToDataTypeString(col.UpdateOptions.DataType)
}
func (col *Column) ToDataTypeString(data DataType) string {
	switch data {
	case SerialType:
		return "SERIAL"
	case StringType:
		if col.Length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.Length)
		}
		return "VARCHAR" // Default to VARCHAR without length if no length specified
	case IntegerType:
		return "INTEGER"
	case DecimalType:
		if col.Precision > 0 && col.Scale > 0 {
			// Handle DECIMAL(precision, scale)
			return fmt.Sprintf("DECIMAL(%d,%d)", col.Precision, col.Scale)
		} else if col.Precision > 0 && col.Scale == 0 {
			// Handle DECIMAL(precision)
			return fmt.Sprintf("DECIMAL(%d)", col.Precision)
		}
		return "DECIMAL"
	case BooleanType:
		return "BOOLEAN"
	case DateType:
		return "DATE"
	case DateTimeType:
		return "TIMESTAMP"
	default:
		return "VARCHAR" // Default to TEXT if type is unknown
	}
}

// IsValidDataType checks if the given data type is valid
func IsValidDataType(dataType DataType) bool {
	switch dataType {
	case SerialType, StringType, IntegerType, DecimalType, BooleanType, DateType, DateTimeType:
		return true
	default:
		return false
	}
}
