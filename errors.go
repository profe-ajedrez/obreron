package obreron

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by builders during Build/BuildInto validation.
//
// All errors are stable contract: callers may rely on errors.Is.
var (
	// Structural validation
	ErrEmptyFrom  = errors.New("obreron: SELECT/DELETE requires a FROM clause")
	ErrEmptyTable = errors.New("obreron: UPDATE/INSERT requires a target table")
	ErrMissingSet = errors.New("obreron: UPDATE requires at least one SET clause")
	ErrNoCols     = errors.New("obreron: INSERT requires at least one column")

	// Paging validation
	ErrInvalidLimit  = errors.New("obreron: LIMIT must be a non-negative integer")
	ErrInvalidOffset = errors.New("obreron: OFFSET must be a non-negative integer")

	// Placeholder / params validation
	ErrTooManyParams       = errors.New("obreron: too many parameters")
	ErrPlaceholderMismatch = errors.New("obreron: placeholder count does not match args")
	ErrReservedMarker      = errors.New("obreron: SQL fragment contains reserved placeholder marker byte")

	// Dialect feature gating
	ErrUnsupportedByDialect = errors.New("obreron: feature not supported by this dialect")
)

// BuildError wraps a sentinel error with additional context about where
// (Op) and under which dialect (Dialect) the failure occurred.
//
// It supports errors.Is / errors.As via Unwrap.
//
// Example error string:
//
//	obreron [FROM/postgres]: obreron: SELECT/DELETE requires a FROM clause
type BuildError struct {
	Op      string // e.g. "SELECT", "FROM", "WHERE", "INSERT"
	Dialect string // dialect name at the time of failure
	Err     error  // underlying sentinel error
}

func (e *BuildError) Error() string {
	if e == nil {
		return "obreron: <nil>"
	}
	// Keep formatting stable; callers may snapshot logs.
	return fmt.Sprintf("obreron [%s/%s]: %v", e.Op, e.Dialect, e.Err)
}

func (e *BuildError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
