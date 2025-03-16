package gomb_test

import (
	"testing"

	gomb "github.com/nandrechetan/gomb/internal"
)

func TestIndex(t *testing.T) {
	t.Run("Basic Index Creation", func(t *testing.T) {
		idx := gomb.NewIndex("idx_users_email")
		idx.OnTable("users")
		idx.AddColumn("email")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_users_email ON users (email)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Unique Index", func(t *testing.T) {
		idx := gomb.NewIndex("idx_users_username")
		idx.OnTable("users")
		idx.AddColumn("username")
		idx.SetUnique()

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE UNIQUE INDEX idx_users_username ON users (username)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Multi-Column Index", func(t *testing.T) {
		idx := gomb.NewIndex("idx_orders_customer_date")
		idx.OnTable("orders")
		idx.MultiColumnIndex("customer_id", "order_date")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_orders_customer_date ON orders (customer_id, order_date)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Concurrent Index Creation", func(t *testing.T) {
		idx := gomb.NewIndex("idx_products_category")
		idx.OnTable("products")
		idx.AddColumn("category_id")
		idx.SetConcurrently()

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX CONCURRENTLY idx_products_category ON products (category_id)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Partial Index with WHERE Clause", func(t *testing.T) {
		idx := gomb.NewIndex("idx_orders_large")
		idx.OnTable("orders")
		idx.AddColumn("id")
		idx.PartialIndex("total_amount > 1000")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_orders_large ON orders (id) WHERE total_amount > 1000"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Expression Index", func(t *testing.T) {
		idx := gomb.NewIndex("idx_users_lower_email")
		idx.OnTable("users")
		idx.ExpressionIndex("LOWER(email)")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_users_lower_email ON users (LOWER(email))"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("GIN Index for JSON", func(t *testing.T) {
		idx := gomb.NewIndex("idx_products_metadata")
		idx.OnTable("products")
		idx.AddColumn("metadata")
		idx.SetMethod("gin")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_products_metadata ON products USING gin (metadata)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Covering Index with INCLUDE", func(t *testing.T) {
		idx := gomb.NewIndex("idx_orders_customer")
		idx.OnTable("orders")
		idx.AddColumn("customer_id")
		idx.AddIncludeColumn("order_date")
		idx.AddIncludeColumn("status")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_orders_customer ON orders (customer_id) INCLUDE (order_date, status)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Index with Tablespace", func(t *testing.T) {
		idx := gomb.NewIndex("idx_large_table")
		idx.OnTable("large_data")
		idx.AddColumn("id")
		idx.SetTablespace("fast_ssd")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_large_table ON large_data (id) TABLESPACE fast_ssd"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Index with WITH options", func(t *testing.T) {
		idx := gomb.NewIndex("idx_users_perf")
		idx.OnTable("users")
		idx.AddColumn("created_at")
		idx.AddWithOption("fillfactor=70")
		idx.AddWithOption("pages_per_range=4")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_users_perf ON users (created_at) WITH (fillfactor=70, pages_per_range=4)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Schema-qualified Index", func(t *testing.T) {
		idx := gomb.NewIndex("idx_users_email")
		idx.OnTable("users")
		idx.AddColumn("email")
		idx.SetSchema("auth")

		sql, err := idx.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "CREATE INDEX idx_users_email ON auth.users (email)"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		testCases := []struct {
			name     string
			setupIdx func() *gomb.Index
		}{
			{
				name: "Missing name",
				setupIdx: func() *gomb.Index {
					idx := gomb.NewIndex("")
					idx.OnTable("users")
					idx.AddColumn("email")
					return idx
				},
			},
			{
				name: "Missing table",
				setupIdx: func() *gomb.Index {
					idx := gomb.NewIndex("idx_users_email")
					idx.AddColumn("email")
					return idx
				},
			},
			{
				name: "Missing columns",
				setupIdx: func() *gomb.Index {
					idx := gomb.NewIndex("idx_users_email")
					idx.OnTable("users")
					return idx
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				idx := tc.setupIdx()
				_, err := idx.ToSQL()
				if err == nil {
					t.Errorf("Expected error for %s, but got nil", tc.name)
				}
			})
		}
	})
}

func TestDropIndex(t *testing.T) {
	t.Run("Basic Drop Index", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX idx_users_email"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index IF EXISTS", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetIfExists()

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX IF EXISTS idx_users_email"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index CONCURRENTLY", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetConcurrently()

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX CONCURRENTLY idx_users_email"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index CASCADE", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetCascade()

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX idx_users_email CASCADE"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index RESTRICT", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetRestrict()

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX idx_users_email RESTRICT"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index with Schema", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetSchema("auth")

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX auth.idx_users_email"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Complex Drop Index", func(t *testing.T) {
		drop := gomb.NewDropIndex("idx_users_email")
		drop.SetIfExists().SetConcurrently().SetCascade().SetSchema("auth")

		sql, err := drop.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "DROP INDEX CONCURRENTLY IF EXISTS auth.idx_users_email CASCADE"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Drop Index Missing Name", func(t *testing.T) {
		drop := gomb.NewDropIndex("")
		_, err := drop.ToSQL()
		if err == nil {
			t.Errorf("Expected error for missing index name, but got nil")
		}
	})
}

func TestRenameIndex(t *testing.T) {
	t.Run("Basic Rename Index", func(t *testing.T) {
		rename := gomb.NewRenameIndex("idx_old", "idx_new")

		sql, err := rename.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "ALTER INDEX idx_old RENAME TO idx_new"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Rename Index with Schema", func(t *testing.T) {
		rename := gomb.NewRenameIndex("idx_old", "idx_new")
		rename.SetSchema("auth")

		sql, err := rename.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "ALTER INDEX auth.idx_old RENAME TO idx_new"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})

	t.Run("Rename Index Missing Names", func(t *testing.T) {
		testCases := []struct {
			name      string
			oldName   string
			newName   string
			expectErr bool
		}{
			{
				name:      "Both names missing",
				oldName:   "",
				newName:   "",
				expectErr: true,
			},
			{
				name:      "Old name missing",
				oldName:   "",
				newName:   "idx_new",
				expectErr: true,
			},
			{
				name:      "New name missing",
				oldName:   "idx_old",
				newName:   "",
				expectErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				rename := gomb.NewRenameIndex(tc.oldName, tc.newName)
				_, err := rename.ToSQL()
				if tc.expectErr && err == nil {
					t.Errorf("Expected error for %s, but got nil", tc.name)
				} else if !tc.expectErr && err != nil {
					t.Errorf("Unexpected error for %s: %v", tc.name, err)
				}
			})
		}
	})
}

func TestReindex(t *testing.T) {
	t.Run("Reindex Table", func(t *testing.T) {
		reindex := gomb.NewReindex("TABLE", "users")

		sql, err := reindex.ToSQL()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "REINDEX TABLE users"
		if sql != expected {
			t.Errorf("Expected SQL: %s, got: %s", expected, sql)
		}
	})
}
