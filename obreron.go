package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

// DB is the factory entrypoint for creating statement builders.
//
// NOTE: Builders are NOT concurrency-safe. One goroutine should use one builder instance.
// See the project spec for details.
type DB struct {
	d dialect.Dialect
}

// New constructs a DB factory using the provided SQL dialect.
func New(d dialect.Dialect) *DB { return &DB{d: d} }

// Select starts a SELECT statement builder.
func (db *DB) Select() *SelectStm { return newSelect(db.d) }

// Insert starts an INSERT statement builder.
func (db *DB) Insert() *InsertStatement { return newInsert(db.d) }

// Update starts an UPDATE statement builder.
func (db *DB) Update(table string) *UpdateStm { return newUpdate(db.d, table) }

// Delete starts a DELETE statement builder.
func (db *DB) Delete() *DeleteStm { return newDelete(db.d) }

// Select starts a SELECT statement builder without using a DB factory.
func Select(d dialect.Dialect) *SelectStm { return newSelect(d) }

// Insert starts an INSERT statement builder without using a DB factory.
func Insert(d dialect.Dialect) *InsertStatement { return newInsert(d) }

// Update starts an UPDATE statement builder without using a DB factory.
func Update(d dialect.Dialect, table string) *UpdateStm { return newUpdate(d, table) }

// Delete starts a DELETE statement builder without using a DB factory.
func Delete(d dialect.Dialect) *DeleteStm { return newDelete(d) }
