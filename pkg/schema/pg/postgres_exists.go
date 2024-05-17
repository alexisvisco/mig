package pg

import (
	"context"
	"fmt"
	"github.com/alexisvisco/amigo/pkg/schema"
	"github.com/georgysavva/scany/v2/dbscan"
)

// ConstraintExist checks if a constraint exists in the Table.
func (p *Schema) ConstraintExist(tableName schema.TableName, constraintName string) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM information_schema.table_constraints WHERE table_name = $1 AND constraint_name = $2 and constraint_schema = $3)"

	row, err := p.db.QueryContext(p.Context.Context, query, tableName.Name(), constraintName, tableName.Schema())
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while checking if constraint exists: %w", err))
		return false
	}

	err = dbscan.ScanOne(&result, row)
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while scanning constraint existence: %w", err))
		return false
	}

	return result
}

// ColumnExist checks if a column exists in the Table.
func (p *Schema) ColumnExist(tableName schema.TableName, columnName string) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2 and table_schema = $3)"

	row, err := p.db.QueryContext(p.Context.Context, query, tableName.Name(), columnName, tableName.Schema())
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while checking if column exists: %w", err))
		return false
	}

	err = dbscan.ScanOne(&result, row)
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while scanning column existence: %w", err))
		return false
	}

	return result
}

// TableExist checks if a table exists in the database.
func (p *Schema) TableExist(tableName schema.TableName) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = $1 AND table_schema = $2)"

	row, err := p.db.QueryContext(p.Context.Context, query, tableName.Name(), tableName.Schema())
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while checking if table exists: %w", err))
		return false
	}

	err = dbscan.ScanOne(&result, row)
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while scanning table existence: %w", err))
		return false
	}

	return result
}

// IndexExist checks if an index exists in the Table.
func (p *Schema) IndexExist(tableName schema.TableName, indexName string) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_indexes WHERE tablename = $1 AND indexname = $2 and schemaname = $3)"

	row, err := p.db.QueryContext(p.Context.Context, query, tableName.Name(), indexName, tableName.Schema())
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while checking if index exists: %w", err))
		return false
	}

	err = dbscan.ScanOne(&result, row)
	if err != nil {
		p.Context.RaiseError(fmt.Errorf("error while scanning index existence: %w", err))
		return false
	}

	return result
}

func (p *Schema) PrimaryKeyExist(tableName schema.TableName) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM information_schema.table_constraints WHERE table_name = $1 AND constraint_type = 'PRIMARY KEY')"

	row, err := p.db.QueryContext(context.Background(), query, tableName.Name())
	if err != nil {
		return false
	}

	err = dbscan.ScanOne(&result, row)
	if err != nil {
		return false
	}

	return result
}
