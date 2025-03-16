package gomb_test

import (
	"testing"

	gomb "github.com/nandrechetan/gomb/internal"
)

func TestMetadataBuilderToSQL(t *testing.T) {
	tests := []struct {
		name        string
		column      *gomb.Column
		expectedSQL string
		expectError bool
	}{
		{
			name: "Test SERIAL column with PrimaryKey",
			column: gomb.NewColumn("id").
				SetDataType(gomb.IntegerType).
				SetPrimaryKey(),
			expectedSQL: "id INTEGER PRIMARY KEY",
			expectError: false,
		},
		{
			name: "Test Integer column with PrimaryKey",
			column: gomb.NewColumn("id").
				SetDataType(gomb.SerialType).
				SetPrimaryKey(),
			expectedSQL: "id SERIAL PRIMARY KEY",
			expectError: false,
		},
		{
			name: "Test String column with UNIQUE constraint",
			column: gomb.NewColumn("username").
				SetDataType(gomb.StringType).
				SetUnique(),
			expectedSQL: "username VARCHAR UNIQUE",
			expectError: false,
		},
		{
			name: "Test String column with NOT NULL",
			column: gomb.NewColumn("email").
				SetDataType(gomb.StringType).
				SetNotNull(),
			expectedSQL: "email VARCHAR NOT NULL",
			expectError: false,
		},
		{
			name: "Test String column with CHECK constraint",
			column: gomb.NewColumn("age").
				SetDataType(gomb.IntegerType).
				SetCheck("(age >= 18)"),
			expectedSQL: "age INTEGER CHECK (age >= 18)",
			expectError: false,
		},
		{
			name: "Test String column with Foreign Key reference",
			column: gomb.NewColumn("user_id").
				SetDataType(gomb.IntegerType).
				SetReferences(gomb.T("users"), gomb.C("id")),
			expectedSQL: "user_id INTEGER REFERENCES users(id)",
			expectError: false,
		},
		{
			name: "Test Column with Generated expression",
			column: gomb.NewColumn("created_at").
				SetDataType(gomb.DateTimeType).
				SetGenerated("CURRENT_TIMESTAMP"),
			expectedSQL: "created_at TIMESTAMP GENERATED ALWAYS AS (CURRENT_TIMESTAMP)",
			expectError: false,
		},
		// New test cases
		{
			name: "Test String column with length and NOT NULL",
			column: gomb.NewColumn("phone").
				SetDataType(gomb.StringType).
				SetLength(15).
				SetNotNull(),
			expectedSQL: "phone VARCHAR(15) NOT NULL",
			expectError: false,
		},
		{
			name: "Test Numeric column with precision and scale",
			column: gomb.NewColumn("price").
				SetDataType(gomb.DecimalType).
				SetPrecision(10).SetScale(2),
			expectedSQL: "price DECIMAL(10,2)",
			expectError: false,
		},
		{
			name: "Test column with multiple constraints",
			column: gomb.NewColumn("email").
				SetDataType(gomb.StringType).
				SetLength(255).
				SetNotNull().
				SetUnique().
				SetCheck("(email LIKE '%@%.%')"),
			expectedSQL: "email VARCHAR(255) NOT NULL UNIQUE CHECK (email LIKE '%@%.%')",
			expectError: false,
		},
		{
			name: "Test Foreign Key with ON DELETE CASCADE",
			column: gomb.NewColumn("user_id").
				SetDataType(gomb.IntegerType).
				SetReferencesOnDeleteCascade(gomb.T("users"), gomb.C("id")),
			expectedSQL: "user_id INTEGER REFERENCES users(id) ON DELETE CASCADE",
			expectError: false,
		},
		{
			name: "Test column with default value",
			column: gomb.NewColumn("status").
				SetDataType(gomb.StringType).
				SetDefault("active"),
			expectedSQL: "status VARCHAR DEFAULT 'active'",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, err := tt.column.ToSQL()
			if tt.expectError && err == nil {
				t.Errorf("Expected error, but got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if sql != tt.expectedSQL {
				t.Errorf("Expected SQL: %s, but got: %s", tt.expectedSQL, sql)
			}
		})
	}

	t.Run("Test Alter Table with New Column", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetDataType(gomb.StringType).SetLength(10).SetReferences(gomb.T("crmuser"), gomb.C("id"))
		table.AddColumn(ownerIdcolumn)
		is_deleteColumn := gomb.NewColumn("is_delete").SetDataType(gomb.BooleanType).SetDefault(gomb.DefaultFalse)
		table.AddColumn(is_deleteColumn)
		orderDateColumn := gomb.NewColumn("order_date").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp)
		table.AddColumn(orderDateColumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account ADD COLUMN ownerId VARCHAR(10) REFERENCES crmuser(id), ADD COLUMN is_delete BOOLEAN DEFAULT FALSE, ADD COLUMN order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})

	t.Run("Test Delete Column with Alter Column", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId")
		table.DropColumn(ownerIdcolumn)
		is_deleteColumn := gomb.NewColumn("is_delete")
		table.DropColumn(is_deleteColumn)
		orderDateColumn := gomb.NewColumn("order_date")
		table.DropColumn(orderDateColumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account DROP COLUMN ownerId, DROP COLUMN is_delete, DROP COLUMN order_date"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})

	t.Run("Test Modify Column with Alter Column Name", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetNewName("owner_id")
		table.AlterColumn(ownerIdcolumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account RENAME COLUMN ownerId TO owner_id"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})

	t.Run("Test Modify Column with Alter Column Data Type", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetDataType(gomb.IntegerType).SetNewDataType(gomb.StringType)
		table.AlterColumn(ownerIdcolumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account ALTER COLUMN INTEGER TYPE VARCHAR"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})
	/*
	   // New test cases for ALTER TABLE operations

	   	t.Run("Test Alter Table with Multiple Operations", func(t *testing.T) {
	   		table := gomb.NewAlterTable("Products")

	   		// Add multiple columns with constraints
	   		priceColumn := gomb.NewColumn("price").SetDataType(gomb.NumericType).SetPrecision(10, 2).SetNotNull()
	   		table.AddColumn(priceColumn)

	   		descColumn := gomb.NewColumn("description").SetDataType(gomb.TextType)
	   		table.AddColumn(descColumn)

	   		// Drop a column
	   		oldColumn := gomb.NewColumn("old_field")
	   		table.DropColumn(oldColumn)

	   		// Rename a column
	   		renameColumn := gomb.NewColumn("created").SetNewName("created_at")
	   		table.AlterColumn(renameColumn)

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "ALTER TABLE Products ADD COLUMN price NUMERIC(10,2) NOT NULL, ADD COLUMN description TEXT, DROP COLUMN old_field, RENAME COLUMN created TO created_at"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Alter Table with Column Default Modification", func(t *testing.T) {
	   		table := gomb.NewAlterTable("Users")
	   		statusColumn := gomb.NewColumn("status").SetDefault("'active'")
	   		table.AlterColumnDefault(statusColumn)

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "ALTER TABLE Users ALTER COLUMN status SET DEFAULT 'active'"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Alter Table with Adding Constraints", func(t *testing.T) {
	   		table := gomb.NewAlterTable("Orders")
	   		table.AddConstraint("fk_customer", "FOREIGN KEY (customer_id) REFERENCES customers(id)")
	   		table.AddConstraint("uq_order_number", "UNIQUE (order_number)")

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "ALTER TABLE Orders ADD CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id), ADD CONSTRAINT uq_order_number UNIQUE (order_number)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Alter Table with Dropping Constraints", func(t *testing.T) {
	   		table := gomb.NewAlterTable("Inventory")
	   		table.DropConstraint("fk_product")
	   		table.DropConstraint("check_qty_positive")

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "ALTER TABLE Inventory DROP CONSTRAINT fk_product, DROP CONSTRAINT check_qty_positive"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Index", func(t *testing.T) {
	   		idx := gomb.NewIndex("idx_accounts_name").
	   			OnTable("accounts").
	   			AddColumn("name").
	   			SetUnique()

	   		genratedSQL, err := idx.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE UNIQUE INDEX idx_accounts_name ON accounts (name)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Multi-Column Index", func(t *testing.T) {
	   		idx := gomb.NewIndex("idx_orders_customer_date").
	   			OnTable("orders").
	   			AddColumn("customer_id").
	   			AddColumn("order_date").
	   			SetConcurrently()

	   		genratedSQL, err := idx.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE INDEX CONCURRENTLY idx_orders_customer_date ON orders (customer_id, order_date)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Table", func(t *testing.T) {
	   		table := gomb.NewCreateTable("employees")

	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		nameCol := gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100).SetNotNull()
	   		emailCol := gomb.NewColumn("email").SetDataType(gomb.StringType).SetLength(255).SetUnique()
	   		deptCol := gomb.NewColumn("department_id").SetDataType(gomb.IntegerType).SetReferences(gomb.T("departments"), gomb.C("id"))
	   		createdCol := gomb.NewColumn("created_at").SetDataType(gomb.TimestampType).SetDefault(gomb.DefaultCurrentTimestamp)

	   		table.AddColumn(idCol)
	   		table.AddColumn(nameCol)
	   		table.AddColumn(emailCol)
	   		table.AddColumn(deptCol)
	   		table.AddColumn(createdCol)

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE employees (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, email VARCHAR(255) UNIQUE, department_id INTEGER REFERENCES departments(id), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Table with Table Constraints", func(t *testing.T) {
	   		table := gomb.NewCreateTable("orders")

	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		custIdCol := gomb.NewColumn("customer_id").SetDataType(gomb.IntegerType).SetNotNull()
	   		orderNumCol := gomb.NewColumn("order_number").SetDataType(gomb.StringType).SetLength(20).SetNotNull()
	   		totalCol := gomb.NewColumn("total_amount").SetDataType(gomb.NumericType).SetPrecision(10, 2).SetNotNull()
	   		createdCol := gomb.NewColumn("created_at").SetDataType(gomb.TimestampType).SetDefault(gomb.DefaultCurrentTimestamp)

	   		table.AddColumn(idCol)
	   		table.AddColumn(custIdCol)
	   		table.AddColumn(orderNumCol)
	   		table.AddColumn(totalCol)
	   		table.AddColumn(createdCol)

	   		table.AddConstraint("fk_customer", "FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE")
	   		table.AddConstraint("uq_order_number", "UNIQUE (order_number)")
	   		table.AddConstraint("check_total", "CHECK (total_amount > 0)")

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE orders (id SERIAL PRIMARY KEY, customer_id INTEGER NOT NULL, order_number VARCHAR(20) NOT NULL, total_amount NUMERIC(10,2) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE, CONSTRAINT uq_order_number UNIQUE (order_number), CONSTRAINT check_total CHECK (total_amount > 0))"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Drop Table", func(t *testing.T) {
	   		dropTable := gomb.NewDropTable("temp_logs").SetIfExists().SetCascade()

	   		genratedSQL, err := dropTable.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "DROP TABLE IF EXISTS temp_logs CASCADE"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Complex Table Creation with Composite Primary Key", func(t *testing.T) {
	   		table := gomb.NewCreateTable("order_items")

	   		orderIdCol := gomb.NewColumn("order_id").SetDataType(gomb.IntegerType).SetNotNull()
	   		productIdCol := gomb.NewColumn("product_id").SetDataType(gomb.IntegerType).SetNotNull()
	   		qtyCol := gomb.NewColumn("quantity").SetDataType(gomb.IntegerType).SetNotNull().SetDefault("1").SetCheck("(quantity > 0)")
	   		priceCol := gomb.NewColumn("unit_price").SetDataType(gomb.NumericType).SetPrecision(10, 2).SetNotNull()

	   		table.AddColumn(orderIdCol)
	   		table.AddColumn(productIdCol)
	   		table.AddColumn(qtyCol)
	   		table.AddColumn(priceCol)

	   		table.AddConstraint("pk_order_items", "PRIMARY KEY (order_id, product_id)")
	   		table.AddConstraint("fk_order", "FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE")
	   		table.AddConstraint("fk_product", "FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT")

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE order_items (order_id INTEGER NOT NULL, product_id INTEGER NOT NULL, quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0), unit_price NUMERIC(10,2) NOT NULL, CONSTRAINT pk_order_items PRIMARY KEY (order_id, product_id), CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE, CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Table with Inheritance", func(t *testing.T) {
	   		table := gomb.NewCreateTable("customer_logs")

	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		messageCol := gomb.NewColumn("message").SetDataType(gomb.TextType).SetNotNull()

	   		table.AddColumn(idCol)
	   		table.AddColumn(messageCol)
	   		table.SetInheritsFrom("logs")

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE customer_logs (id SERIAL PRIMARY KEY, message TEXT NOT NULL) INHERITS (logs)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Batch Execution", func(t *testing.T) {
	   		batch := gomb.NewBatch()

	   		// Create table
	   		createTable := gomb.NewCreateTable("products")
	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		nameCol := gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100).SetNotNull()
	   		createTable.AddColumn(idCol)
	   		createTable.AddColumn(nameCol)
	   		batch.AddStatement(createTable)

	   		// Create index
	   		idx := gomb.NewIndex("idx_products_name").OnTable("products").AddColumn("name")
	   		batch.AddStatement(idx)

	   		// Insert data
	   		insert := gomb.NewInsert("products").
	   			Columns("name").
	   			Values("'Test Product 1'").
	   			Values("'Test Product 2'")
	   		batch.AddStatement(insert)

	   		genratedSQL, err := batch.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE products (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL);\nCREATE INDEX idx_products_name ON products (name);\nINSERT INTO products (name) VALUES ('Test Product 1'), ('Test Product 2');"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Table With Schema", func(t *testing.T) {
	   		table := gomb.NewCreateTable("users").SetSchema("auth")

	   		idCol := gomb.NewColumn("id").SetDataType(gomb.UuidType).SetDefault("gen_random_uuid()").SetPrimaryKey()
	   		usernameCol := gomb.NewColumn("username").SetDataType(gomb.StringType).SetLength(50).SetNotNull().SetUnique()
	   		passwordCol := gomb.NewColumn("password_hash").SetDataType(gomb.StringType).SetLength(100).SetNotNull()

	   		table.AddColumn(idCol)
	   		table.AddColumn(usernameCol)
	   		table.AddColumn(passwordCol)

	   		genratedSQL, err := table.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "CREATE TABLE auth.users (id UUID DEFAULT gen_random_uuid() PRIMARY KEY, username VARCHAR(50) NOT NULL UNIQUE, password_hash VARCHAR(100) NOT NULL)"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Transaction", func(t *testing.T) {
	   		txn := gomb.NewTransaction()

	   		// Create table
	   		createTable := gomb.NewCreateTable("audit_logs")
	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		actionCol := gomb.NewColumn("action").SetDataType(gomb.StringType).SetLength(50).SetNotNull()
	   		createTable.AddColumn(idCol)
	   		createTable.AddColumn(actionCol)
	   		txn.AddStatement(createTable)

	   		// Insert data
	   		insert := gomb.NewInsert("audit_logs").
	   			Columns("action").
	   			Values("'Initial setup'")
	   		txn.AddStatement(insert)

	   		genratedSQL, err := txn.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "BEGIN;\nCREATE TABLE audit_logs (id SERIAL PRIMARY KEY, action VARCHAR(50) NOT NULL);\nINSERT INTO audit_logs (action) VALUES ('Initial setup');\nCOMMIT;"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Insert with Returning", func(t *testing.T) {
	   		insert := gomb.NewInsert("users").
	   			Columns("username", "email", "active").
	   			Values("'john_doe'", "'john@example.com'", "true").
	   			Returning("id", "created_at")

	   		genratedSQL, err := insert.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "INSERT INTO users (username, email, active) VALUES ('john_doe', 'john@example.com', true) RETURNING id, created_at"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Update with Conditions", func(t *testing.T) {
	   		update := gomb.NewUpdate("products").
	   			Set("price", "price * 1.1").
	   			Set("updated_at", "NOW()").
	   			Where("category_id = 5").
	   			AndWhere("price < 100")

	   		genratedSQL, err := update.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "UPDATE products SET price = price * 1.1, updated_at = NOW() WHERE category_id = 5 AND price < 100"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Delete with Conditions", func(t *testing.T) {
	   		delete := gomb.NewDelete("temp_records").
	   			Where("created_at < NOW() - INTERVAL '30 days'")

	   		genratedSQL, err := delete.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedSQL := "DELETE FROM temp_records WHERE created_at < NOW() - INTERVAL '30 days'"

	   		if genratedSQL != expectedSQL {
	   			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
	   		}
	   	})

	   	t.Run("Test Create Migration", func(t *testing.T) {
	   		migration := gomb.NewMigration("create_customer_table")

	   		// Up operations
	   		createTable := gomb.NewCreateTable("customers")
	   		idCol := gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey()
	   		nameCol := gomb.NewColumn("name").SetDataType(gomb.StringType).SetLength(100).SetNotNull()
	   		emailCol := gomb.NewColumn("email").SetDataType(gomb.StringType).SetLength(255).SetUnique()
	   		createTable.AddColumn(idCol)
	   		createTable.AddColumn(nameCol)
	   		createTable.AddColumn(emailCol)
	   		migration.AddUpStatement(createTable)

	   		createIndex := gomb.NewIndex("idx_customers_email").OnTable("customers").AddColumn("email")
	   		migration.AddUpStatement(createIndex)

	   		// Down operations
	   		dropIndex := gomb.NewDropIndex("idx_customers_email")
	   		migration.AddDownStatement(dropIndex)

	   		dropTable := gomb.NewDropTable("customers")
	   		migration.AddDownStatement(dropTable)

	   		upSQL, downSQL, err := migration.ToSQL()
	   		if err != nil {
	   			t.Errorf("Unexpected error: %v", err)
	   		}

	   		expectedUpSQL := "-- Migration: create_customer_table (UP)\nCREATE TABLE customers (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, email VARCHAR(255) UNIQUE);\nCREATE INDEX idx_customers_email ON customers (email);"
	   		expectedDownSQL := "-- Migration: create_customer_table (DOWN)\nDROP INDEX idx_customers_email;\nDROP TABLE customers;"

	   		if upSQL != expectedUpSQL {
	   			t.Errorf("Generated UP SQL mismatch.\nExpected: %s\nGot: %s", expectedUpSQL, upSQL)
	   		}

	   		if downSQL != expectedDownSQL {
	   			t.Errorf("Generated DOWN SQL mismatch.\nExpected: %s\nGot: %s", expectedDownSQL, downSQL)
	   		}
	   	})
	*/
}
func TestColumn_ToSQL(t *testing.T) {
	tests := []struct {
		name        string
		column      *gomb.Column
		expectedSQL string
		expectError bool
	}{
		{
			name: "Test SERIAL column with PrimaryKey",
			column: gomb.NewColumn("id").
				SetDataType(gomb.IntegerType).
				SetPrimaryKey(),
			expectedSQL: "id INTEGER PRIMARY KEY",
			expectError: false,
		},
		{
			name: "Test Integer column with PrimaryKey",
			column: gomb.NewColumn("id").
				SetDataType(gomb.SerialType).
				SetPrimaryKey(),
			expectedSQL: "id SERIAL PRIMARY KEY",
			expectError: false,
		},
		{
			name: "Test String column with UNIQUE constraint",
			column: gomb.NewColumn("username").
				SetDataType(gomb.StringType).
				SetUnique(),
			expectedSQL: "username VARCHAR UNIQUE",
			expectError: false,
		},
		{
			name: "Test String column with NOT NULL",
			column: gomb.NewColumn("email").
				SetDataType(gomb.StringType).
				SetNotNull(),
			expectedSQL: "email VARCHAR NOT NULL",
			expectError: false,
		},
		{
			name: "Test String column with CHECK constraint",
			column: gomb.NewColumn("age").
				SetDataType(gomb.IntegerType).
				SetCheck("(age >= 18)"),
			expectedSQL: "age INTEGER CHECK (age >= 18)",
			expectError: false,
		},
		{
			name: "Test String column with Foreign Key reference",
			column: gomb.NewColumn("user_id").
				SetDataType(gomb.IntegerType).
				SetReferences(gomb.T("users"), gomb.C("id")),
			expectedSQL: "user_id INTEGER REFERENCES users(id)",
			expectError: false,
		},
		{
			name: "Test Column with Generated expression",
			column: gomb.NewColumn("created_at").
				SetDataType(gomb.DateTimeType).
				SetGenerated("CURRENT_TIMESTAMP"),
			expectedSQL: "created_at TIMESTAMP GENERATED ALWAYS AS (CURRENT_TIMESTAMP)",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, err := tt.column.ToSQL()
			if tt.expectError && err == nil {
				t.Errorf("Expected error, but got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if sql != tt.expectedSQL {
				t.Errorf("Expected SQL: %s, but got: %s", tt.expectedSQL, sql)
			}
		})
	}

	t.Run("Test Alter Table with New Column", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetDataType(gomb.StringType).SetLength(10).SetReferences(gomb.T("crmuser"), gomb.C("id"))
		table.AddColumn(ownerIdcolumn)
		is_deleteColumn := gomb.NewColumn("is_delete").SetDataType(gomb.BooleanType).SetDefault(gomb.DefaultFalse)
		table.AddColumn(is_deleteColumn)
		orderDateColumn := gomb.NewColumn("order_date").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp)
		table.AddColumn(orderDateColumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account ADD COLUMN ownerId VARCHAR(10) REFERENCES crmuser(id), ADD COLUMN is_delete BOOLEAN DEFAULT FALSE, ADD COLUMN order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})

	t.Run("Test Delete Column with Alter Column", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId")
		table.DropColumn(ownerIdcolumn)
		is_deleteColumn := gomb.NewColumn("is_delete")
		table.DropColumn(is_deleteColumn)
		orderDateColumn := gomb.NewColumn("order_date")
		table.DropColumn(orderDateColumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account DROP COLUMN ownerId, DROP COLUMN is_delete, DROP COLUMN order_date"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})

	t.Run("Test Modify Column with Alter Column Name", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetNewName("owner_id")
		table.AlterColumn(ownerIdcolumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account RENAME COLUMN ownerId TO owner_id"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})
	t.Run("Test Modify Column with Alter Column Data Type", func(t *testing.T) {
		table := gomb.NewAlterTable("Account")
		ownerIdcolumn := gomb.NewColumn("ownerId").SetDataType(gomb.IntegerType).SetNewDataType(gomb.StringType)
		table.AlterColumn(ownerIdcolumn)
		genratedSQL, err := table.ToSQL()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedSQL := "ALTER TABLE Account ALTER COLUMN INTEGER TYPE VARCHAR"

		if genratedSQL != expectedSQL {
			t.Errorf("Generated SQL mismatch.\nExpected: %s\nGot: %s", expectedSQL, genratedSQL)
		}
	})
}
