package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

// DB is the factory entrypoint for statement builders.
//
// Builders created via DB are NOT concurrency-safe. One goroutine must
// own one builder. See the project spec for the recommended [sync.Pool] pattern.
type DB struct {
	d dialect.Dialect
}

// New constructs a DB factory for the given dialect.
func New(d dialect.Dialect) *DB { return &DB{d: d} }

// Select starts a SELECT statement builder.
func (db *DB) Select() *SelectStm { return newSelect(db.d) }

// Insert starts an INSERT statement builder.
func (db *DB) Insert() *InsertStatement { return newInsert(db.d) }

// Update starts an UPDATE statement builder.
func (db *DB) Update() *UpdateStm { return newUpdate(db.d) }

// Delete starts a DELETE statement builder.
func (db *DB) Delete() *DeleteStm { return newDelete(db.d) }

// Package-level convenience constructors — no need for a DB instance.

// Select starts a SELECT statement builder.
func Select(d dialect.Dialect) *SelectStm { return newSelect(d) }

// Insert starts an INSERT statement builder.
func Insert(d dialect.Dialect) *InsertStatement { return newInsert(d) }

// Update starts an UPDATE statement builder.
func Update(d dialect.Dialect) *UpdateStm { return newUpdate(d) }

// Delete starts a DELETE statement builder.
func Delete(d dialect.Dialect) *DeleteStm { return newDelete(d) }
