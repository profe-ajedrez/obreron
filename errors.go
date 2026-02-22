package obreron

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyFrom            = errors.New("obreron: SELECT/DELETE requires a FROM clause")
	ErrEmptyTable           = errors.New("obreron: UPDATE/INSERT requires a target table")
	ErrMissingSet           = errors.New("obreron: UPDATE requires at least one SET clause")
	ErrNoCols               = errors.New("obreron: INSERT requires at least one column")
	ErrInvalidLimit         = errors.New("obreron: LIMIT must be a non-negative integer")
	ErrTooManyParams        = errors.New("obreron: parameter count exceeds dialect maximum")
	ErrPlaceholderMismatch  = errors.New("obreron: placeholder count does not match argument count")
	ErrUnsupportedByDialect = errors.New("obreron: feature not supported by this dialect")
)

type BuildError struct {
	Op      string // método donde ocurrió: "FROM", "WHERE", "RETURNING", etc.
	Dialect string // Dialect.Name() en el momento del error
	Err     error  // sentinel subyacente — compatible con errors.Is / errors.As
}

func (e *BuildError) Error() string {
	return fmt.Sprintf("obreron [%s/%s]: %v", e.Op, e.Dialect, e.Err)
}

func (e *BuildError) Unwrap() error { return e.Err }
