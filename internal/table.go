package gomb

import (
	"fmt"
	"strings"
)

// NewAlterTable initializes and returns a new AlterTable instance
func NewAlterTable(name string) *AlterTable {
	return &AlterTable{
		TableName:  name,
		Operations: make([]ColumnOperation, 0),
	}
}

// AddColumn adds a new column to the table
func (t *AlterTable) AddColumn(column *Column) *AlterTable {
	if column != nil {
		t.Operations = append(t.Operations, ColumnOperation{
			Operation: AddColumnOp,
			Column:    column,
		})
	}
	return t
}

// DropColumn marks a column for deletion
func (t *AlterTable) DropColumn(column *Column) *AlterTable {
	if column != nil {
		t.Operations = append(t.Operations, ColumnOperation{
			Operation: DropColumnOp,
			Column:    column,
		})
	}
	return t
}

// AlterColumn marks a column for alteration
func (t *AlterTable) AlterColumn(column *Column) *AlterTable {
	if column != nil {
		// Determine if it's a rename or type change
		if column.NewName != "" {
			t.Operations = append(t.Operations, ColumnOperation{
				Operation: RenameColumnOp,
				Column:    column,
			})
		} else if column.NewDataType != "" {
			t.Operations = append(t.Operations, ColumnOperation{
				Operation: AlterColumnTypeOp,
				Column:    column,
			})
		}
	}
	return t
}

// ToSQL generates the SQL statement for ALTER TABLE
func (t *AlterTable) ToSQL() (string, []error) {
	errors := t.Validate()
	if len(errors) > 0 {
		return "", errors
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s", t.TableName))

	// Process operations
	operationDefs := make([]string, 0, len(t.Operations))
	for _, op := range t.Operations {
		switch op.Operation {
		case AddColumnOp:
			colSQL, err := op.Column.ToSQL()
			if err != nil {
				errors = append(errors, err)
				continue
			}
			operationDefs = append(operationDefs, "ADD COLUMN "+colSQL)
		case DropColumnOp:
			operationDefs = append(operationDefs, "DROP COLUMN "+op.Column.Name)
		case RenameColumnOp:
			operationDefs = append(operationDefs, "RENAME COLUMN "+op.Column.Name+" TO "+op.Column.NewName)
		case AlterColumnTypeOp:
			operationDefs = append(operationDefs, "ALTER COLUMN "+op.Column.ToDataType()+" TYPE "+op.Column.ToNewDataType())
		}
	}

	if len(operationDefs) == 0 {
		errors = append(errors, fmt.Errorf("no valid operations defined for table %s", t.TableName))
		return "", errors
	}

	builder.WriteString(" " + strings.Join(operationDefs, ", "))

	// Add table-level comment if provided
	sql := builder.String()
	if t.Comment != "" {
		sql = sql + fmt.Sprintf(" COMMENT ON TABLE %s IS '%s'", t.TableName, t.Comment)
	}

	return sql, nil
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

// Validate validates the alter table operation
func (t *AlterTable) Validate() []error {
	var errors []error

	// Check if table name is empty
	if t.TableName == "" {
		errors = append(errors, fmt.Errorf("table name cannot be empty"))
	}

	// Check if there are operations
	if len(t.Operations) == 0 {
		errors = append(errors, fmt.Errorf("alter table must have at least one operation"))
	}

	return errors
}
