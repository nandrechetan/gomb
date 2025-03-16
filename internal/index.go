package gomb

import (
	"fmt"
	"strings"
)

// NewIndex creates a new index builder
func NewIndex(name string) *Index {
	return &Index{
		name:    name,
		columns: []string{},
		using:   "btree", // Default to btree
	}
}

// OnTable sets the table for the index
func (idx *Index) OnTable(table string) *Index {
	idx.table = table
	return idx
}

// AddColumn adds a column to the index
func (idx *Index) AddColumn(column string) *Index {
	idx.columns = append(idx.columns, column)
	return idx
}

// SetUnique marks the index as unique
func (idx *Index) SetUnique() *Index {
	idx.unique = true
	return idx
}

// SetConcurrently sets the index to be built concurrently
func (idx *Index) SetConcurrently() *Index {
	idx.concurrently = true
	return idx
}

// SetMethod sets the index method (btree, hash, gist, gin, etc.)
func (idx *Index) SetMethod(method string) *Index {
	idx.method = method
	return idx
}

// SetWhere adds a where clause to the index
func (idx *Index) SetWhere(condition string) *Index {
	idx.where = condition
	return idx
}

// SetSchema sets the schema for the index
func (idx *Index) SetSchema(schema string) *Index {
	idx.schema = schema
	return idx
}

// AddIncludeColumn adds a column to the INCLUDE clause (for covering indexes)
func (idx *Index) AddIncludeColumn(column string) *Index {
	idx.includeColumns = append(idx.includeColumns, column)
	return idx
}

// SetTablespace sets the tablespace for the index
func (idx *Index) SetTablespace(tablespace string) *Index {
	idx.tablespace = tablespace
	return idx
}

// AddWithOption adds an option to the WITH clause
func (idx *Index) AddWithOption(option string) *Index {
	idx.withOptions = append(idx.withOptions, option)
	return idx
}

// ToSQL generates the SQL for creating the index
func (idx *Index) ToSQL() (string, error) {
	if idx.name == "" {
		return "", fmt.Errorf("index name is required")
	}

	if idx.table == "" {
		return "", fmt.Errorf("table name is required")
	}

	if len(idx.columns) == 0 {
		return "", fmt.Errorf("at least one column is required for an index")
	}

	var sql strings.Builder

	sql.WriteString("CREATE ")

	if idx.unique {
		sql.WriteString("UNIQUE ")
	}

	sql.WriteString("INDEX ")

	if idx.concurrently {
		sql.WriteString("CONCURRENTLY ")
	}

	sql.WriteString(idx.name)
	sql.WriteString(" ON ")

	if idx.schema != "" {
		sql.WriteString(idx.schema)
		sql.WriteString(".")
	}

	sql.WriteString(idx.table)

	if idx.method != "" {
		sql.WriteString(" USING ")
		sql.WriteString(idx.method)
	}

	sql.WriteString(" (")
	sql.WriteString(strings.Join(idx.columns, ", "))
	sql.WriteString(")")

	if len(idx.includeColumns) > 0 {
		sql.WriteString(" INCLUDE (")
		sql.WriteString(strings.Join(idx.includeColumns, ", "))
		sql.WriteString(")")
	}

	if idx.where != "" {
		sql.WriteString(" WHERE ")
		sql.WriteString(idx.where)
	}

	if len(idx.withOptions) > 0 {
		sql.WriteString(" WITH (")
		sql.WriteString(strings.Join(idx.withOptions, ", "))
		sql.WriteString(")")
	}

	if idx.tablespace != "" {
		sql.WriteString(" TABLESPACE ")
		sql.WriteString(idx.tablespace)
	}

	return sql.String(), nil
}

// DropIndex represents a DROP INDEX operation
type DropIndex struct {
	name         string
	ifExists     bool
	concurrently bool
	cascade      bool
	restrict     bool
	schema       string
}

// NewDropIndex creates a new drop index builder
func NewDropIndex(name string) *DropIndex {
	return &DropIndex{
		name: name,
	}
}

// SetIfExists adds IF EXISTS to the drop statement
func (di *DropIndex) SetIfExists() *DropIndex {
	di.ifExists = true
	return di
}

// SetConcurrently sets the drop to be done concurrently
func (di *DropIndex) SetConcurrently() *DropIndex {
	di.concurrently = true
	return di
}

// SetCascade adds CASCADE to the drop statement
func (di *DropIndex) SetCascade() *DropIndex {
	di.cascade = true
	di.restrict = false
	return di
}

// SetRestrict adds RESTRICT to the drop statement
func (di *DropIndex) SetRestrict() *DropIndex {
	di.restrict = true
	di.cascade = false
	return di
}

// SetSchema sets the schema for the index
func (di *DropIndex) SetSchema(schema string) *DropIndex {
	di.schema = schema
	return di
}

// ToSQL generates the SQL for dropping the index
func (di *DropIndex) ToSQL() (string, error) {
	if di.name == "" {
		return "", fmt.Errorf("index name is required")
	}

	var sql strings.Builder

	sql.WriteString("DROP INDEX ")

	if di.concurrently {
		sql.WriteString("CONCURRENTLY ")
	}

	if di.ifExists {
		sql.WriteString("IF EXISTS ")
	}

	if di.schema != "" {
		sql.WriteString(di.schema)
		sql.WriteString(".")
	}

	sql.WriteString(di.name)

	if di.cascade {
		sql.WriteString(" CASCADE")
	} else if di.restrict {
		sql.WriteString(" RESTRICT")
	}

	return sql.String(), nil
}

// RenameIndex represents a RENAME INDEX operation
type RenameIndex struct {
	oldName string
	newName string
	schema  string
}

// NewRenameIndex creates a new rename index builder
func NewRenameIndex(oldName, newName string) *RenameIndex {
	return &RenameIndex{
		oldName: oldName,
		newName: newName,
	}
}

// SetSchema sets the schema for the index
func (ri *RenameIndex) SetSchema(schema string) *RenameIndex {
	ri.schema = schema
	return ri
}

// ToSQL generates the SQL for renaming the index
func (ri *RenameIndex) ToSQL() (string, error) {
	if ri.oldName == "" || ri.newName == "" {
		return "", fmt.Errorf("both old and new index names are required")
	}

	var sql strings.Builder

	sql.WriteString("ALTER INDEX ")

	if ri.schema != "" {
		sql.WriteString(ri.schema)
		sql.WriteString(".")
	}

	sql.WriteString(ri.oldName)
	sql.WriteString(" RENAME TO ")
	sql.WriteString(ri.newName)

	return sql.String(), nil
}

// ReindexOperation represents a REINDEX operation
type ReindexOperation struct {
	target       string // INDEX, TABLE, SCHEMA, DATABASE, SYSTEM
	name         string
	concurrently bool
}

// NewReindex creates a new reindex builder
func NewReindex(target, name string) *ReindexOperation {
	return &ReindexOperation{
		target: strings.ToUpper(target),
		name:   name,
	}
}

// SetConcurrently sets the reindex to be done concurrently
func (ro *ReindexOperation) SetConcurrently() *ReindexOperation {
	ro.concurrently = true
	return ro
}

// ToSQL generates the SQL for the reindex operation
func (ro *ReindexOperation) ToSQL() (string, error) {
	if ro.target == "" {
		return "", fmt.Errorf("reindex target is required")
	}

	if ro.name == "" && ro.target != "SYSTEM" && ro.target != "DATABASE" {
		return "", fmt.Errorf("name is required for REINDEX %s", ro.target)
	}

	var sql strings.Builder

	sql.WriteString("REINDEX ")

	if ro.concurrently {
		sql.WriteString("CONCURRENTLY ")
	}

	sql.WriteString(ro.target)

	if ro.name != "" {
		sql.WriteString(" ")
		sql.WriteString(ro.name)
	}

	return sql.String(), nil
}

// SetIndexTablespace represents an ALTER INDEX SET TABLESPACE operation
type SetIndexTablespace struct {
	indexName  string
	tablespace string
	nowait     bool
	schema     string
}

// NewSetIndexTablespace creates a new set index tablespace builder
func NewSetIndexTablespace(indexName, tablespace string) *SetIndexTablespace {
	return &SetIndexTablespace{
		indexName:  indexName,
		tablespace: tablespace,
	}
}

// SetNowait adds NOWAIT to the statement
func (sit *SetIndexTablespace) SetNowait() *SetIndexTablespace {
	sit.nowait = true
	return sit
}

// SetSchema sets the schema for the index
func (sit *SetIndexTablespace) SetSchema(schema string) *SetIndexTablespace {
	sit.schema = schema
	return sit
}

// ToSQL generates the SQL for the set tablespace operation
func (sit *SetIndexTablespace) ToSQL() (string, error) {
	if sit.indexName == "" {
		return "", fmt.Errorf("index name is required")
	}

	if sit.tablespace == "" {
		return "", fmt.Errorf("tablespace name is required")
	}

	var sql strings.Builder

	sql.WriteString("ALTER INDEX ")

	if sit.schema != "" {
		sql.WriteString(sit.schema)
		sql.WriteString(".")
	}

	sql.WriteString(sit.indexName)
	sql.WriteString(" SET TABLESPACE ")
	sql.WriteString(sit.tablespace)

	if sit.nowait {
		sql.WriteString(" NOWAIT")
	}

	return sql.String(), nil
}

// PartialIndex creates a new partial index with a WHERE clause
func (idx *Index) PartialIndex(condition string) *Index {
	return idx.SetWhere(condition)
}

// ExpressionIndex adds an expression to the index instead of a simple column
func (idx *Index) ExpressionIndex(expression string) *Index {
	idx.columns = append(idx.columns, expression)
	return idx
}

// MultiColumnIndex adds multiple columns to the index at once
func (idx *Index) MultiColumnIndex(columns ...string) *Index {
	idx.columns = append(idx.columns, columns...)
	return idx
}

// IsStatement implementation for SQL generation interface
func (idx *Index) IsStatement() {}

// IsStatement implementation for SQL generation interface
func (di *DropIndex) IsStatement() {}

// IsStatement implementation for SQL generation interface
func (ri *RenameIndex) IsStatement() {}

// IsStatement implementation for SQL generation interface
func (ro *ReindexOperation) IsStatement() {}

// IsStatement implementation for SQL generation interface
func (sit *SetIndexTablespace) IsStatement() {}
