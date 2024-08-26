package pg

import "github.com/jackc/pgx/v5"

// RowsAdapter makes pgx.Rows compliant with the dbsql.Rows interface.
// See dbsql.Rows for details.
type RowsAdapter struct {
	pgx.Rows
}

// NewRowsAdapter returns a new rowsAdapter instance.
func NewRowsAdapter(rows pgx.Rows) *RowsAdapter {
	return &RowsAdapter{Rows: rows}
}

// Columns implements the dbscan.Rows.Columns method.
func (ra RowsAdapter) Columns() ([]string, error) {
	columns := make([]string, len(ra.Rows.FieldDescriptions()))
	for i, fd := range ra.Rows.FieldDescriptions() {
		columns[i] = fd.Name
	}
	return columns, nil
}

// Close implements the dbscan.Rows.Close method.
func (ra RowsAdapter) Close() error {
	ra.Rows.Close()
	return nil
}

// NextResultSet is currently always returning false.
func (ra RowsAdapter) NextResultSet() bool {
	// TODO: when pgx issue #308 and #1512 and is fixed maybe we can do something here.
	return false
}
