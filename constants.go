package gomb

// Define a custom type for data types
type DataType string
type DefaultValue string
type Constraint string

// AlterTableOperation represents the type of operation to perform
type AlterTableOperation int

const (
	AddColumnOp AlterTableOperation = iota
	DropColumnOp
	RenameColumnOp
	AlterColumnTypeOp
)

// Define constants for each data type as a custom type
const (
	SerialType   DataType = "serial"
	StringType   DataType = "string"
	IntegerType  DataType = "integer"
	DecimalType  DataType = "decimal"
	BooleanType  DataType = "boolean"
	DateType     DataType = "date"
	DateTimeType DataType = "datetime"
)

// Constants for PostgreSQL data types (prefix 'Pg' for PostgreSQL)
const (
	DefaultNull DefaultValue = "NULL"
	// Boolean defaults
	DefaultTrue  DefaultValue = "TRUE"
	DefaultFalse DefaultValue = "FALSE"

	// Date/Time defaults
	DefaultCurrentTimestamp DefaultValue = "CURRENT_TIMESTAMP"
	DefaultCurrentDate      DefaultValue = "CURRENT_DATE"
	DefaultCurrentTime      DefaultValue = "CURRENT_TIME"
	DefaultLocalTime        DefaultValue = "LOCALTIME"
	DefaultLocalTimestamp   DefaultValue = "LOCALTIMESTAMP"

	// Constraints
	PrimaryKey Constraint = "PRIMARY KEY"
	NotNull    Constraint = "NOT NULL"
	Unique     Constraint = "UNIQUE"
)
