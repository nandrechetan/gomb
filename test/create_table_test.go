package gomb_test

import (
	"testing"

	gomb "github.com/nandrechetan/gomb/internal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTable_ToSQL(t *testing.T) {
	tests := []struct {
		name       string
		table      *gomb.Table
		wantSQL    string
		wantErrors bool
	}{
		{
			name: "Basic Table",
			table: func() *gomb.Table {
				table := gomb.NewTable("users")
				table.AddColumn(gomb.NewColumn("id").SetPrimaryKey().SetDataType(gomb.SerialType))
				table.AddColumn(gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(50))
				return table
			}(),
			wantSQL:    "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(50))",
			wantErrors: false,
		},
		{
			name: "Table With All Column Types",
			table: func() *gomb.Table {
				table := gomb.NewTable("all_types")
				table.AddColumn(gomb.NewColumn("id").SetPrimaryKey().SetDataType(gomb.SerialType))
				table.AddColumn(gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100).SetNotNull())
				table.AddColumn(gomb.NewColumn("description").SetDataType(gomb.StringType))
				table.AddColumn(gomb.NewColumn("active").SetDataType(gomb.BooleanType).SetDefault(gomb.DefaultTrue))
				table.AddColumn(gomb.NewColumn("count").SetDataType(gomb.IntegerType))
				table.AddColumn(gomb.NewColumn("price").SetDataType(gomb.DecimalType).SetPrecision(10).SetScale(2))
				table.AddColumn(gomb.NewColumn("created_at").SetDataType(gomb.DateTimeType).SetDefault("CURRENT_TIMESTAMP"))
				table.AddColumn(gomb.NewColumn("updated_at").SetDataType(gomb.DateTimeType))
				// table.AddColumn(gomb.NewColumn("json_data").SetDataType(gomb.JSONType))
				return table
			}(),
			wantSQL:    "CREATE TABLE all_types (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, description VARCHAR, active BOOLEAN DEFAULT TRUE, count INTEGER, price DECIMAL(10,2), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP)",
			wantErrors: false,
		},
		{
			name: "Table With Comment",
			table: func() *gomb.Table {
				table := gomb.NewTable("products")
				table.Comment = "Products table stores all product information"
				table.AddColumn(gomb.NewColumn("id").SetPrimaryKey().SetDataType(gomb.SerialType))
				table.AddColumn(gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100))
				return table
			}(),
			wantSQL:    "CREATE TABLE products (id SERIAL PRIMARY KEY, name VARCHAR(100)) COMMENT ON TABLE products IS 'Products table stores all product information'",
			wantErrors: false,
		},
		{
			name: "Columns With Comments",
			table: func() *gomb.Table {
				table := gomb.NewTable("employees")
				idCol := gomb.NewColumn("id").SetPrimaryKey().SetDataType(gomb.SerialType)
				idCol.Comment = "Primary identifier for employees"
				table.AddColumn(idCol)

				nameCol := gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100).SetNotNull()
				nameCol.Comment = "Employee full name"
				table.AddColumn(nameCol)
				return table
			}(),
			wantSQL:    "CREATE TABLE employees (id SERIAL PRIMARY KEY COMMENT 'Primary identifier for employees', name VARCHAR(100) NOT NULL COMMENT 'Employee full name')",
			wantErrors: false,
		},
		{
			name:       "Empty Table Name",
			table:      gomb.NewTable(""),
			wantSQL:    "",
			wantErrors: true,
		},
		{
			name: "Table Without Columns",
			table: func() *gomb.Table {
				return gomb.NewTable("empty_table")
			}(),
			wantSQL:    "",
			wantErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, errors := tt.table.ToSQL()

			if tt.wantErrors {
				assert.NotEmpty(t, errors, "Expected errors but got none")
			} else {
				assert.Empty(t, errors, "Expected no errors but got: %v", errors)
				assert.Equal(t, tt.wantSQL, sql, "SQL doesn't match expected")
			}
		})
	}
}

func TestComplex_Scenarios(t *testing.T) {
	t.Run("Complete Users Table", func(t *testing.T) {
		table := gomb.NewTable("users")
		table.Comment = "Store user information"

		// Add columns
		table.AddColumn(gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey().SetComment("User ID"))
		table.AddColumn(gomb.NewColumn("username").SetDataType(gomb.StringType).SetLength(50).SetNotNull().SetComment("Unique username"))
		table.AddColumn(gomb.NewColumn("email").SetDataType(gomb.StringType).SetLength(255).SetNotNull().SetComment("User email address"))
		table.AddColumn(gomb.NewColumn("password_hash").SetDataType(gomb.StringType).SetLength(100).SetNotNull())
		table.AddColumn(gomb.NewColumn("first_name").SetDataType(gomb.StringType).SetLength(50))
		table.AddColumn(gomb.NewColumn("last_name").SetDataType(gomb.StringType).SetLength(50))
		table.AddColumn(gomb.NewColumn("birth_date").SetDataType(gomb.DateType))
		table.AddColumn(gomb.NewColumn("is_active").SetDataType(gomb.BooleanType).SetDefault(gomb.DefaultTrue))
		table.AddColumn(gomb.NewColumn("login_count").SetDataType(gomb.IntegerType).SetDefault(0))
		table.AddColumn(gomb.NewColumn("last_login").SetDataType(gomb.DateTimeType))
		table.AddColumn(gomb.NewColumn("created_at").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp))
		table.AddColumn(gomb.NewColumn("updated_at").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp))

		sql, errors := table.ToSQL()
		assert.Empty(t, errors, "Expected no errors but got: %v", errors)
		expectedSQL := "CREATE TABLE users (id SERIAL PRIMARY KEY COMMENT 'User ID', username VARCHAR(50) NOT NULL COMMENT 'Unique username', email VARCHAR(255) NOT NULL COMMENT 'User email address', password_hash VARCHAR(100) NOT NULL, first_name VARCHAR(50), last_name VARCHAR(50), birth_date DATE, is_active BOOLEAN DEFAULT TRUE, login_count INTEGER DEFAULT 0, last_login TIMESTAMP, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP) COMMENT ON TABLE users IS 'Store user information'"
		assert.Equal(t, expectedSQL, sql)
	})
}
