package gomb_test

import (
	"testing"

	gomb "github.com/nandrechetan/gomb/internal"
	"github.com/stretchr/testify/assert"
)

func TestAlterTable_ToSQL(t *testing.T) {
	tests := []struct {
		name       string
		alterTable *gomb.AlterTable
		wantSQL    string
		wantErrors bool
	}{
		{
			name: "Add Column",
			alterTable: func() *gomb.AlterTable {
				alter := gomb.NewAlterTable("users")
				alter.AddColumn(gomb.NewColumn("email").SetDataType(gomb.StringType).SetLength(255).SetNotNull())
				return alter
			}(),
			wantSQL:    "ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL",
			wantErrors: false,
		},
		{
			name: "Drop Column",
			alterTable: func() *gomb.AlterTable {
				alter := gomb.NewAlterTable("users")
				alter.DropColumn(gomb.NewColumn("temp_field"))
				return alter
			}(),
			wantSQL:    "ALTER TABLE users DROP COLUMN temp_field",
			wantErrors: false,
		},
		{
			name: "Rename Column",
			alterTable: func() *gomb.AlterTable {
				alter := gomb.NewAlterTable("users")
				col := gomb.NewColumn("username")
				col.SetNewName("login_name")
				alter.AlterColumn(col)
				return alter
			}(),
			wantSQL:    "ALTER TABLE users RENAME COLUMN username TO login_name",
			wantErrors: false,
		},
		{
			name: "Multiple Operations",
			alterTable: func() *gomb.AlterTable {
				alter := gomb.NewAlterTable("products")
				alter.AddColumn(gomb.NewColumn("category_id").SetDataType(gomb.IntegerType).SetNotNull())
				alter.DropColumn(gomb.NewColumn("old_category"))

				renameCol := gomb.NewColumn("desc")
				renameCol.SetNewName("description")
				alter.AlterColumn(renameCol)

				return alter
			}(),
			wantSQL:    "ALTER TABLE products ADD COLUMN category_id INTEGER NOT NULL, DROP COLUMN old_category, RENAME COLUMN desc TO description",
			wantErrors: false,
		},
		{
			name: "Alter Table With Comment",
			alterTable: func() *gomb.AlterTable {
				alter := gomb.NewAlterTable("orders")
				alter.Comment = "Updated orders table"
				alter.AddColumn(gomb.NewColumn("status").SetDataType(gomb.StringType).SetLength(20))
				return alter
			}(),
			wantSQL:    "ALTER TABLE orders ADD COLUMN status VARCHAR(20) COMMENT ON TABLE orders IS 'Updated orders table'",
			wantErrors: false,
		},
		{
			name:       "Empty Table Name",
			alterTable: gomb.NewAlterTable(""),
			wantSQL:    "",
			wantErrors: true,
		},
		{
			name: "No Operations",
			alterTable: func() *gomb.AlterTable {
				return gomb.NewAlterTable("table_without_ops")
			}(),
			wantSQL:    "",
			wantErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, errors := tt.alterTable.ToSQL()

			if tt.wantErrors {
				assert.NotEmpty(t, errors, "Expected errors but got none")
			} else {
				assert.Empty(t, errors, "Expected no errors but got: %v", errors)
				assert.Equal(t, tt.wantSQL, sql, "SQL doesn't match expected")
			}
		})
	}
}
